package com.quentinguidee.dao

import com.quentinguidee.models.OAuthMethod
import com.quentinguidee.models.OAuthMethodPrivate
import com.quentinguidee.models.OAuthMethods
import org.jetbrains.exposed.sql.ResultRow
import org.jetbrains.exposed.sql.selectAll

class OAuthMethodsDAO {
    private fun toOAuthMethodPrivate(row: ResultRow) = OAuthMethodPrivate(
        name = row[OAuthMethods.name],
        color = row[OAuthMethods.color],
        clientID = row[OAuthMethods.clientID],
        clientSecret = row[OAuthMethods.clientSecret],
    )

    private fun toOAuthMethod(row: ResultRow) = OAuthMethod(
        name = row[OAuthMethods.name],
        color = row[OAuthMethods.color],
    )

    fun getAll() = OAuthMethods
        .slice(OAuthMethods.name, OAuthMethods.color)
        .selectAll()
        .map(::toOAuthMethod)
}

val oAuthMethodsDAO = OAuthMethodsDAO()
