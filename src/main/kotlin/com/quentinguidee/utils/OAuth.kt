package com.quentinguidee.utils

import com.quentinguidee.client
import io.ktor.client.call.*
import io.ktor.client.request.*
import io.ktor.http.*
import kotlinx.serialization.SerialName
import kotlinx.serialization.Serializable
import java.util.*

data class OAuthConfig(
    val clientID: String,
    val clientSecret: String,
    val authorizeURL: String,
    val accessTokenURL: String,
    val redirectURL: String,
)

@Serializable
data class OAuthExchangeResponse(
    @SerialName("access_token")
    val accessToken: String,
)

class OAuth(private val config: OAuthConfig) {
    fun getLoginURL(): String {
        val url = URLBuilder(config.authorizeURL)

        url.parameters.set("client_id", config.clientID)
        url.parameters.set("redirect_uri", config.redirectURL)
        url.parameters.set("state", UUID.randomUUID().toString())

        return url.buildString()
    }

    suspend fun exchange(config: OAuthConfig, code: String): OAuthExchangeResponse = client
        .request(config.accessTokenURL) {
            method = HttpMethod.Post
            parameter("client_id", config.clientID)
            parameter("client_secret", config.clientSecret)
            parameter("code", code)
            parameter("redirect_uri", config.redirectURL)
            headers {
                accept(ContentType.Application.Json)
            }
        }
        .body()
}
