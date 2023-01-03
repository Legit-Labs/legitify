#!/usr/bin/env python3

import sys
import yaml
import os

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
### Remediation
{remediation_string}
"""

    tmp = f"""## {title}
policy name: {policy_name}

severity: {severity}

### Description
{description}
"""
    if len(threat) > 0:
        threat_string = "".join([f"{line}\n" for index, line in enumerate(threat)])
        tmp = f"""{tmp}
### Threat Example(s)
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

def print_usage():
    print("python gen-gh-pages-docs.py docs_file output_directory")

if __name__ == '__main__':
    if len(sys.argv) != 3:
        print_usage()
        exit(1)
    create_policy_docs(sys.argv[1], sys.argv[2])
