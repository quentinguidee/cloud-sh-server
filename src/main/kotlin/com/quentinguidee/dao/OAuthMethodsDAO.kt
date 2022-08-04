package com.quentinguidee.dao

import com.quentinguidee.models.OAuthMethod
import com.quentinguidee.models.OAuthMethodPrivate
import com.quentinguidee.models.OAuthMethods
import org.jetbrains.exposed.sql.ResultRow
import org.jetbrains.exposed.sql.insert
import org.jetbrains.exposed.sql.select
import org.jetbrains.exposed.sql.selectAll

class OAuthMethodsDAO {
    private fun toOAuthMethod(row: ResultRow) = OAuthMethod(
        id = row[OAuthMethods.id].value,
        name = row[OAuthMethods.name],
        displayName = row[OAuthMethods.displayName],
        color = row[OAuthMethods.color],
        clientID = row[OAuthMethods.clientID],
        authorizeURL = row[OAuthMethods.authorizeURL],
        accessTokenURL = row[OAuthMethods.accessTokenURL],
        redirectURL = row[OAuthMethods.redirectURL],
    )

    private fun toOAuthMethodPrivate(row: ResultRow) = OAuthMethodPrivate(
        id = row[OAuthMethods.id].value,
        name = row[OAuthMethods.name],
        displayName = row[OAuthMethods.displayName],
        color = row[OAuthMethods.color],
        clientID = row[OAuthMethods.clientID],
        clientSecret = row[OAuthMethods.clientSecret],
        authorizeURL = row[OAuthMethods.authorizeURL],
        accessTokenURL = row[OAuthMethods.accessTokenURL],
        redirectURL = row[OAuthMethods.redirectURL],
    )

    fun create(method: OAuthMethodPrivate) =
        create(
            method.name,
            method.displayName,
            method.color,
            method.clientID,
            method.clientSecret,
            method.authorizeURL,
            method.accessTokenURL,
            method.redirectURL,
        )


    fun create(
        name: String,
        displayName: String,
        color: String,
        clientID: String,
        clientSecret: String,
        authorizeURL: String,
        accessTokenURL: String,
        redirectURL: String,
    ) = OAuthMethods.insert {
        it[OAuthMethods.name] = name
        it[OAuthMethods.displayName] = displayName
        it[OAuthMethods.color] = color
        it[OAuthMethods.clientID] = clientID
        it[OAuthMethods.clientSecret] = clientSecret
        it[OAuthMethods.authorizeURL] = authorizeURL
        it[OAuthMethods.accessTokenURL] = accessTokenURL
        it[OAuthMethods.redirectURL] = redirectURL
    }

    fun getAll() = OAuthMethods
        .selectAll()
        .map(::toOAuthMethod)

    fun getAllPrivate() = OAuthMethods
        .selectAll()
        .map(::toOAuthMethodPrivate)

    fun get(name: String) = OAuthMethods
        .select { OAuthMethods.name eq name }
        .map(::toOAuthMethod)
        .first()

    fun getPrivate(name: String) = OAuthMethods
        .select { OAuthMethods.name eq name }
        .map(::toOAuthMethodPrivate)
        .first()
}

val oAuthMethodsDAO = OAuthMethodsDAO()
