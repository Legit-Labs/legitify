#!/usr/bin/env python3

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

def create_policy_page(policy, output_dir, parent, grand_parent):
    file_path = os.path.join(output_dir, f"{policy['policy_name']}.md")
    md = gen_policy_markdown(policy)
    title=policy['title']
    final =f"""---
layout: default
title: {title}
parent: {parent}
grand_parent: {grand_parent}
---

{md}
"""
    with open(file_path, 'w') as f:
        f.write(final)

def create_ns_policies(output_dir, ns, docs_yaml, parent):
    ns_dir = os.path.join(output_dir)
    os.mkdir(ns_dir)
    title = f"{ns.title()} Policies"

    file_path = os.path.join(ns_dir, f"index.md")
    file_header=f"""---
layout: default
title: {title}
parent: {parent}
has_children: true
---
"""

    with open(file_path, 'w') as f:
        f.write(file_header)

    for policy in docs_yaml[ns]:
        create_policy_page(policy, ns_dir, title, parent)

    return ns_dir

def create_scm_policy_docs(scm, docs_yaml, output_dir):
    scm_outdir = os.path.join(output_dir, scm)
    os.mkdir(scm_outdir)
    file_path = os.path.join(scm_outdir, f"index.md")
    title = f"{scm_to_pretty_name(scm)} Policies"
    file_header=f"""---
layout: default
title: {title}
has_children: true
---
"""

    with open(file_path, 'w') as f:
        f.write(file_header)

    for ns in docs_yaml:
        store_at = os.path.join(scm_outdir, ns)
        create_ns_policies(store_at, ns, docs_yaml, title)

def create_policy_docs(docs_file, output_dir):
    docs_yaml = get_docs_yaml(docs_file)

    for scm in docs_yaml:
        create_scm_policy_docs(scm, docs_yaml[scm], output_dir)

def create_monomarkdown(docs_file, output_dir):
    table_of_contents = f"{format_header(HeaderSize.H1, 'Table of contents')}\n"

    docs_yaml = get_docs_yaml(docs_file)
    result = ""

    i = 1
    for scm in docs_yaml:
        pretty_name = scm_to_pretty_name(scm)
        result += f"{format_header(HeaderSize.H1, pretty_name)}\n"
        table_of_contents += f"{i}. [{pretty_name}](#{scm})\n"
        j = 1
        for ns in docs_yaml[scm]:
            table_of_contents += f"\t{j}. [{ns}](#{ns})\n"
            result += f"{format_header(HeaderSize.H2, ns)}\n"
            k = 1
            for policy in docs_yaml[scm][ns]:
                table_of_contents += f"\t\t{k}. [{policy['title']}](#{policy['title'].lower().replace(' ', '-')})\n"
                policy_md = gen_policy_markdown(policy)
                result += f"{policy_md}\n"
                k += 1
            j += 1
        i += 1

    with open(os.path.join(output_dir, "monomarkdown.md"), "w") as f:
        f.write(table_of_contents)
        f.write("\n")
        f.write(result)



if __name__ == "__main__":
    parser = argparse.ArgumentParser()
    parser.add_argument("docs_file", type=str, help="input file name")
    parser.add_argument("output_dir", type=str, help="output directory")
    parser.add_argument("--monomarkdown", action='store_true')

    args = parser.parse_args()

    if not args.monomarkdown:
        create_policy_docs(args.docs_file, args.output_dir)
    else:
        create_monomarkdown(args.docs_file, args.output_dir)
