import argparse
import os
import subprocess
import requests

GITHUB_ORG = "Legit-Labs"
GITHUB_REPO = "homebrew-legit-labs"
GITHB_ORG_AND_REPO=f"{GITHUB_ORG}/{GITHUB_REPO}"
FORMULA_FILE_PATH="legitify.rb"


def create_local_changes(version, arm_filename, arm_sha256, intel_filename, intel_sha256, formula_file_path=FORMULA_FILE_PATH):
    BREW_FORMULA = f"""
class Legitify < Formula
desc "Legitify - open source scm scanning tool by Legit Security"
homepage "https://github.com/Legit-Labs/legitify"

on_macos do
    if Hardware::CPU.intel?
    url "https://github.com/Legit-Labs/legitify/releases/download/{version}/{intel_filename}"
    sha256 "{intel_sha256}"
    version "{version}"
    end 
    if Hardware::CPU.arm? && !Hardware::CPU.is_64_bit?
    url "https://github.com/Legit-Labs/legitify/releases/download/{version}/{arm_filename}"
    sha256 "{arm_sha256}"
    version "{version}" 
    end
end

    def install
    bin.install "legitify"
    end
end
"""
    with open(formula_file_path, 'w') as f:
        f.write(BREW_FORMULA)

def checkout_new_branch(bump_version):
    sanitized_version = bump_version.replace('.','_')
    branch_name = f'feat/update_{sanitized_version}'
    process = subprocess.Popen(['git', 'checkout', '-b', branch_name],
                     stderr=subprocess.PIPE)
    stderr = process.communicate()
    return branch_name


def commit_and_push():
    process = subprocess.Popen(['git', 'add', FORMULA_FILE_PATH],
                     stderr=subprocess.PIPE)
    stderr = process.communicate()

    process = subprocess.Popen(['git', 'commit', '-m', 'Bump brew formula'],
                     stderr=subprocess.PIPE)
    stderr = process.communicate()

    push_repo()
    
def push_repo():
    GITHUB_USER = os.environ['GITHUB_USER']
    HOMEBREW_TAP_GITHUB_TOKEN = os.environ['HOMEBREW_TAP_GITHUB_TOKEN']
    process = subprocess.Popen(['git', 'push', f'https://{GITHUB_USER}:{HOMEBREW_TAP_GITHUB_TOKEN}@github.com/{GITHB_ORG_AND_REPO}'],
                     stderr=subprocess.PIPE)
    stderr = process.communicate()


def clone_repo():
    os.chdir('/tmp') # to avoid cloning inside repo
    GITHUB_USER = os.environ['GITHUB_USER']
    HOMEBREW_TAP_GITHUB_TOKEN = os.environ['HOMEBREW_TAP_GITHUB_TOKEN']
    process = subprocess.Popen(['git', 'clone', f'https://{GITHUB_USER}:{HOMEBREW_TAP_GITHUB_TOKEN}@github.com/{GITHB_ORG_AND_REPO}'],
                     stderr=subprocess.PIPE)
    stderr = process.communicate()
    os.chdir(GITHUB_REPO)

    
def create_pull_request(bump_version, head_branch, repo_path=GITHB_ORG_AND_REPO):
    url = f"https://api.github.com/repos/{repo_path}/pulls"
    HOMEBREW_TAP_GITHUB_TOKEN = os.environ['HOMEBREW_TAP_GITHUB_TOKEN']
    headers = {
        "Accept": "application/vnd.github+json",
        "Authorization": f"Bearer {HOMEBREW_TAP_GITHUB_TOKEN}",
        "X-GitHub-Api-Version": "2022-11-28",
    }

    data = {
        "title": f"legitify {bump_version}",
        "body": "Auto generated using github actions and internal scripts",
        "head": head_branch,
        "base": "main",
    }

    response = requests.post(url, headers=headers, json=data)

    if response.status_code == 201:
        print("Pull request created successfully!")
    else:
        print(f"Failed to create pull request. Status code: {response.status_code}")
        print(response.text)
        exit(1)

if __name__ == "__main__":
    parser = argparse.ArgumentParser()
    parser.add_argument("version", type=str, help="release version")
    parser.add_argument("arm_filename", type=str, help="arm file name")
    parser.add_argument("arm_sha256", type=str, help="arm file sha")
    parser.add_argument("intel_filename", type=str, help="intel file name")
    parser.add_argument("intel_sha256", type=str, help="intel file sha")

    args = parser.parse_args()
    clone_repo()
    branch_name = checkout_new_branch(args.version)
    create_local_changes(args.version, args.arm_filename, args.arm_sha256, args.intel_filename, args.intel_sha256)
    commit_and_push()
    create_pull_request(args.version, branch_name)