package com.quentinguidee.dao

import com.quentinguidee.models.OAuthMethodPrivate
import com.quentinguidee.models.OAuthUser
import com.quentinguidee.models.OAuthUsers
import org.jetbrains.exposed.sql.ResultRow
import org.jetbrains.exposed.sql.insert
import org.jetbrains.exposed.sql.select

class OAuthUsersDAO {
    private fun toOAuthUser(row: ResultRow) = OAuthUser(
        userID = row[OAuthUsers.user].value,
        username = row[OAuthUsers.username],
        oAuthMethodID = row[OAuthUsers.oAuthMethod].value
    )

    fun get(username: String) = OAuthUsers
        .select { OAuthUsers.username eq username }
        .map(::toOAuthUser)
        .first()

    fun create(userID: Int, username: String, oAuthMethod: OAuthMethodPrivate) = OAuthUsers
        .insert {
            it[OAuthUsers.user] = userID
            it[OAuthUsers.username] = username
            it[OAuthUsers.oAuthMethod] = oAuthMethod.id
        }.resultedValues!!.map(::toOAuthUser).first()
}

val oAuthUsersDAO = OAuthUsersDAO()
