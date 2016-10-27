-- #####################################################################
-- server.server is a table that contains server configurations.
-- #####################################################################

-- server.tunnel_addr is a listening address to accept client tunnel request.
server.tunnel_addr = "0.0.0.0:9690"

-- server.https_addr is a public address for HTTP connections
server.http_addr = "0.0.0.0:9680"

-- server.https_addr is a public address for HTTPS connections
-- server.https_addr = "0.0.0.0:9681"

-- server.domain is a domain where the tunnels are hosted
server.domain = "xgrok-example.com"

-- server.tls_crt Path to a TLS certificate file
-- server.tls_crt = "path/to/crt",

-- server.tls_key Path to a TLS key file
-- server.tls_key = "path/to/key",

-- server.disable_tcp is to disable TCP protocol proxy.
server.disable_tcp = true


-- #####################################################################
-- user_auth is a table that contains user auth configurations.
-- #####################################################################

user_auth.enable = false
user_auth.tokens = {
    "aaa123abc",
    "bbb456def",
}

-- #####################################################################
-- hooks is a table that contains hooks configurations.
-- #####################################################################

-- hooks.pre_register_tunnel = function()
--
-- end
--
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
