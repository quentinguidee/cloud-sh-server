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
        name = row[OAuthMethods.name],
        displayName = row[OAuthMethods.displayName],
        color = row[OAuthMethods.color],
        clientID = row[OAuthMethods.clientID],
        authorizeURL = row[OAuthMethods.authorizeURL],
        accessTokenURL = row[OAuthMethods.accessTokenURL],
        redirectURL = row[OAuthMethods.redirectURL],
    )

    private fun toOAuthMethodPrivate(row: ResultRow) = OAuthMethodPrivate(
        name = row[OAuthMethods.name],
        displayName = row[OAuthMethods.displayName],
        color = row[OAuthMethods.color],
        clientID = row[OAuthMethods.clientID],
        clientSecret = row[OAuthMethods.clientSecret],
        authorizeURL = row[OAuthMethods.authorizeURL],
        accessTokenURL = row[OAuthMethods.accessTokenURL],
        redirectURL = row[OAuthMethods.redirectURL],
    )

    fun create(oAuthMethod: OAuthMethodPrivate) = OAuthMethods.insert {
        it[name] = oAuthMethod.name
        it[displayName] = oAuthMethod.displayName
        it[color] = oAuthMethod.color
        it[clientID] = oAuthMethod.clientID
        it[clientSecret] = oAuthMethod.clientSecret
        it[authorizeURL] = oAuthMethod.authorizeURL
        it[accessTokenURL] = oAuthMethod.accessTokenURL
        it[redirectURL] = oAuthMethod.redirectURL
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
