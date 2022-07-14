package common.webhooks

has_key(x, k) {
	_ = x[k]
}

is_secure_hook(hook) {
    has_key(hook.config, "insecure_ssl")
    hook.config.insecure_ssl == "0"
    has_key(hook.config, "secret")
}