-- #####################################################################
-- server.server is a table that contains server configurations.
-- #####################################################################

-- server.tunnel_addr is a listening address to accept client tunnel requests.
server.tunnel_addr = "0.0.0.0:9690"

-- server.status_addr is a listening address to provide status API.
server.status_addr = "0.0.0.0:9691"

-- server.https_addr is a public address for HTTP connections
server.http_addr = "0.0.0.0:9680"

-- server.https_addr is a public address for HTTPS connections
-- ex) server.https_addr = "0.0.0.0:443"
server.https_addr = ""

-- server.domain is a domain where the tunnels are hosted
server.domain = "your_xgrok_domain.com"

-- server.tls_crt Path to a TLS certificate file
-- ex) server.tls_crt = "/etc/xgrok/tls/crt"
server.tls_crt = ""

-- server.tls_key Path to a TLS key file
-- ex) server.tls_crt = "/etc/xgrok/tls/key"
server.tls_key = ""

-- server.disable_tcp is to disable TCP protocol proxy.
server.disable_tcp = false

-- server.disable_hostname is to disable using custom hostname proxy.
server.disable_hostname = false


-- #####################################################################
-- user_auth is a table that contains user auth configurations.
-- #####################################################################

-- user_auth.enable activates built-in token base user authentication.
user_auth.enable = false
-- user_auth.tokens is a list of user tokens.
user_auth.tokens = {
    "aaa123abc",
    "bbb456def",
}


-- #####################################################################
-- hooks is a table that contains hook functions.
-- #####################################################################

--
-- Hooks description:
--   In all hook functions, You can return value that represents a error.
--

-- hooks.msg_auth_filter = function(msg_auth)
--
-- end

-- hooks.msg_auth_resp_filter = function(msg_auth_resp)
--
-- end

-- hooks.msg_req_tunnel_filter = function(msg_req_tunnel)
--
-- end

-- hooks.msg_new_tunnel_filter = function(msg_new_tunnel)
--     msg_new_tunnel.custom_props = {
--         { key = "value" },
--         { foo = "bar" },
--     }
-- end

-- hooks.post_register_tunnel = function(tunnel)
--     print("register_tunnel: " .. tunnel.url)
-- end
--
-- hooks.pre_shutdown_tunnel = function(tunnel)
--
-- end
--
-- hooks.post_shutdown_tunnel = function(tunnel)
--
-- end
