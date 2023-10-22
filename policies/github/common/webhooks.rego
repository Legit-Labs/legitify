package common.webhooks

import future.keywords.in

ssl_enabled(hook) {
    "insecure_ssl" in object.keys(hook.config)
    hook.config.insecure_ssl == "0"
}

has_secret(hook) {
    "secret" in object.keys(hook.config)
}
