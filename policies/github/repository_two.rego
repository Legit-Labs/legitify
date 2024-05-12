package repository

import data.common.webhooks as webhookUtils

# METADATA
# scope: rule
# title: Repository Has Stale Secrets
# description: Some of the repository secrets have not been updated for over a year. It is recommended to refresh secret values regularly in order to minimize the risk of breach in case of an information leak.
# custom:
#   requiredEnrichers: [secretsList]
#   remediationSteps:
#      - Enter your repository's landing page
#      - Go to the settings tab
#      - Under the 'Security' title on the left, choose 'Secrets and variables'
#      - Click 'Actions'
#      - Sort secrets by 'Last Updated'
#      - Regenerate every secret older than one year and add the new value to GitHub's secret manager
#   severity: MEDIUM
#   requiredScopes: [repo]
#   threat: Sensitive data may have been inadvertently made public in the past, and an attacker who holds this data may gain access to your current CI and services. In addition, there may be old or unnecessary tokens that have not been inspected and can be used to access sensitive information.
repository_secret_is_stale[stale] := true{
    some index
    secret := input.repository_secrets[index]
    is_stale(secret.updated_at)
    stale := {
        "name" : secret.name,
        "update date" : time.format(secret.updated_at),
    }

}

is_stale(date) {
    ns_per_year := 365 * 24 * 60 * 60 * 1000000000
    diff_ns := time.now_ns() - date
    diff_ns > ns_per_year
}