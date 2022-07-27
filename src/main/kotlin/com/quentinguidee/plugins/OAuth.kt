package com.quentinguidee.plugins

import com.quentinguidee.client
import io.ktor.http.*
import io.ktor.server.application.*
import io.ktor.server.auth.*

fun Application.configureOAuth() {
    val environment = environment
    val httpClient = client

    install(Authentication) {
        oauth("oauth-github") {
            urlProvider = { "http://localhost:8080/auth/github/callback" }
            providerLookup = {
                OAuthServerSettings.OAuth2ServerSettings(
                    name = "github",
                    authorizeUrl = "https://github.com/login/oauth/authorize",
                    accessTokenUrl = "https://github.com/login/oauth/access_token",
                    requestMethod = HttpMethod.Post,
                    clientId = environment.config.property("auth.github.client_id").getString(),
                    clientSecret = environment.config.property("auth.github.client_secret").getString(),
                )
            }
            client = httpClient
        }
    }
}
