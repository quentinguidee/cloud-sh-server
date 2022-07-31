package com.quentinguidee.models

import kotlinx.serialization.SerialName
import kotlinx.serialization.Serializable
import org.jetbrains.exposed.dao.id.IntIdTable

object OAuthMethods : IntIdTable("oauth_methods") {
    val name = varchar("name", 63)
    val color = varchar("color", 7)
    val clientID = varchar("client_id", 63)
    val clientSecret = varchar("client_secret", 63)
}

@Serializable
data class OAuthMethod(
    val name: String,
    val color: String,
)

@Serializable
data class OAuthMethodPrivate(
    val name: String,
    val color: String,
    @SerialName("client_id")
    val clientID: String,
    @SerialName("client_secret")
    val clientSecret: String,
)
