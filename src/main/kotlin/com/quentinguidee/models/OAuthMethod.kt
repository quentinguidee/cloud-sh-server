package com.quentinguidee.models

import com.quentinguidee.client
import io.ktor.client.call.*
import io.ktor.client.request.*
import io.ktor.http.*
import kotlinx.serialization.SerialName
import kotlinx.serialization.Serializable
import org.jetbrains.exposed.dao.id.IntIdTable
import java.util.*

object OAuthMethods : IntIdTable("oauth_methods") {
    val name = varchar("name", 63)
    val displayName = varchar("display_name", 63)
    val color = varchar("color", 7)
    val clientID = varchar("client_id", 63)
    val clientSecret = varchar("client_secret", 63)
    val authorizeURL = varchar("authorize_url", 255)
    val accessTokenURL = varchar("access_token_url", 255)
    val redirectURL = varchar("redirect_url", 255)
}

@Serializable
data class OAuthMethod(
    val id: Int,
    val name: String,
    val displayName: String,
    val color: String,
    val clientID: String,
    val authorizeURL: String,
    val accessTokenURL: String,
    val redirectURL: String,
) {
    fun getLoginURL(): String {
        val url = URLBuilder(authorizeURL)
        url.parameters["client_id"] = clientID
        url.parameters["redirect_uri"] = redirectURL
        url.parameters["state"] = UUID.randomUUID().toString()
        return url.buildString()
    }
}

@Serializable
data class OAuthMethodPrivate(
    val id: Int,
    val name: String,
    val displayName: String,
    val color: String,
    val clientID: String,
    var clientSecret: String,
    val authorizeURL: String,
    val accessTokenURL: String,
    val redirectURL: String,
) {
    suspend fun exchange(code: String): OAuthExchangeResponse = client
        .request(accessTokenURL) {
            method = HttpMethod.Post
            parameter("client_id", clientID)
            parameter("client_secret", clientSecret)
            parameter("code", code)
            parameter("redirect_uri", redirectURL)
            headers {
                accept(ContentType.Application.Json)
            }
        }
        .body()
}

@Serializable
data class OAuthExchangeResponse(
    @SerialName("access_token")
    val accessToken: String,
)
