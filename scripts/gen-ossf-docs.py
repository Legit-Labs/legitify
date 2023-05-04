#!/usr/bin/env python3
import json
import sys
import yaml
import os
import argparse
from enum import Enum


class HeaderSize(Enum):
    H1 = 1
    H2 = 2
    H3 = 3
    H4 = 4
    H5 = 5
    H6 = 6


def format_header(header_size: HeaderSize, text: str) -> str:
    if not isinstance(header_size, HeaderSize):
        raise ValueError("Invalid header size")

    if not isinstance(text, str):
        raise ValueError("Invalid text")

    return f"{'#' * header_size.value} {text}"


def rephrase_ns_name(ns):
    mapping = {
        'organization': 'Organizational Management',
        'group': 'Server',
        'enterprise': 'Server',
        'actions': 'Continuous Integration / Continuous Deployment',
        'runner_group': 'Continuous Integration / Continuous Deployment',
        'repository': 'Repository configuration',
        'project': 'Repository configuration',
        'member': 'Members, Access Control and Permissions'
    }

    return mapping[ns]

def scm_to_pretty_name(scm):
    if scm == 'github': return 'GitHub'
    return 'GitLab'


def get_docs_yaml(docs_file):
    with open(docs_file) as f:
        return yaml.load(f, Loader=yaml.FullLoader)


def gen_policy_markdown(policy):
    policy_name = policy['policy_name']
    title = policy['title']
    description = policy['description']
    severity = policy['severity']
    remediation = policy['remediation']
    threat = policy['threat']

    remediation_string = "".join([f"{index+1}. {line}\n" for index, line in enumerate(remediation)])
    remediation = f"""
{format_header(HeaderSize.H3, "Remediation")}
{remediation_string}
"""

    tmp = f"""{format_header(HeaderSize.H2, title)}
policy name: {policy_name}

severity: {severity}

{format_header(HeaderSize.H3, "Description")}
{description}
"""
    if len(threat) > 0:
        threat_string = "".join([f"{line}\n" for index, line in enumerate(threat)])
        tmp = f"""{tmp}
{format_header(HeaderSize.H3, "Threat Example(s)")}
{threat_string}
"""

    return f"""
{tmp}
{remediation}
"""


def create_policy_page(policy, output_dir):
    file_path = os.path.join(output_dir, f"{policy['policy_name']}.md")
    md = gen_policy_markdown(policy)
    with open(file_path, 'w') as f:
        f.write(md)


def create_ns_policies(output_dir, ns, docs_yaml):
    os.mkdir(output_dir)
    ns_dir = os.path.join(output_dir)

    for policy in docs_yaml[ns]:
        create_policy_page(policy, ns_dir)

    return ns_dir


def create_scm_policy_docs(scm, docs_yaml, output_dir):
    scm_outdir = os.path.join(output_dir, scm)
    os.mkdir(scm_outdir)
    file_path = os.path.join(scm_outdir, f"README.md")

    j = 1
    table_of_contents = ''
    for ns in docs_yaml:
        k = 1
        table_of_contents += f"{format_header(HeaderSize.H3, rephrase_ns_name(ns))}\n"
        for policy in docs_yaml[ns]:
            policy_path = os.path.join(ns, f"{policy['policy_name']}.md")
            table_of_contents += f"{k}. [{policy['title']}]({policy_path})\n"
            k += 1
        j += 1

    with open(file_path, 'w') as f:
        f.write(table_of_contents)

    for ns in docs_yaml:
        store_at = os.path.join(scm_outdir, ns)
        create_ns_policies(store_at, ns, docs_yaml)


def create_policy_docs(docs_file, output_dir):
    docs_yaml = get_docs_yaml(docs_file)

    for scm in docs_yaml:
        create_scm_policy_docs(scm, docs_yaml[scm], output_dir)


def process_root_doc(docs_yaml):
    """
        {
            "category": [
                {"policy_title": [{"link": "link", "icon", "icon"}]
            ]
        }
    """
    icons = {
        "github": '<img src="https://user-images.githubusercontent.com/287526/230375178-2f1f8844-5609-4ef3-b9ac-141c20c43406.svg" alt="GitHub" height="20" width="20">',
        "gitlab": '<img src="https://user-images.githubusercontent.com/287526/230376963-ae9b8a47-4a74-4746-bc83-5b34cc520d40.svg" alt="GitLab" height="20" width="20">'
    }

    result = {}
    for scm in docs_yaml:
        for ns in docs_yaml[scm]:
            category = rephrase_ns_name(ns)
            if category not in result:
                result[category] = {}

            for policy in docs_yaml[scm][ns]:
                policy_link = os.path.join(scm, ns, f"{policy['policy_name']}.md")
                link = {
                    "link": policy_link,
                    "icon": icons[scm]
                }
                policy_title = policy['title']
                if policy_title not in result[category]:
                    result[category][policy_title] = []

                result[category][policy_title].append(link)
    return result


def create_root_doc(docs_file, output_dir):
    docs_yaml = get_docs_yaml(docs_file)
    processed = process_root_doc(docs_yaml)
    result = f"""
## Recommendations

Each specific recommendation below is noted to be applicable to either GitHub or GitLab by use of an appropriate icon, which is linked to the detailed best practice definition: <img src="https://user-images.githubusercontent.com/287526/230375178-2f1f8844-5609-4ef3-b9ac-141c20c43406.svg" alt="GitHub" height="20" width="20"> <img src="https://user-images.githubusercontent.com/287526/230376963-ae9b8a47-4a74-4746-bc83-5b34cc520d40.svg" alt="GitLab" height="20" width="20">

For recommendations only applicable to GitHub or GitLab visit one of the following pages:

- [GitHub Recommendations](github/README.md)
- [GitLab Recommendations](gitlab/README.md)
    
    """

    for category, policies in processed.items():
        result += f"{format_header(header_size=HeaderSize.H3, text=category)} \n"
        for policy, links in policies.items():
            result += f"- {policy}"
            for link in links:
                result += f" [{link['icon']}]({link['link']})"
            result += "\n"
        result += "\n"

    with open(os.path.join(output_dir, "README.md"), "w") as f:
        f.write(result)

if __name__ == "__main__":
    parser = argparse.ArgumentParser()
    parser.add_argument("docs_file", type=str, help="input file name")
    parser.add_argument("output_dir", type=str, help="output directory")

    args = parser.parse_args()

    create_policy_docs(args.docs_file, args.output_dir)
    create_root_doc(args.docs_file, args.output_dir)
