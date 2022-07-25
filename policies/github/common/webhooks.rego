package common.webhooks

has_key(x, k) {
	_ = x[k]
}

ssl_enabled(hook) {
    has_key(hook.config, "insecure_ssl")
    hook.config.insecure_ssl == "0"
}

has_secret(hook) {
    has_key(hook.config, "secret")
}