package com.quentinguidee.dao

import com.quentinguidee.models.GitHubUser
import com.quentinguidee.models.GitHubUsers
import org.jetbrains.exposed.sql.ResultRow
import org.jetbrains.exposed.sql.insert
import org.jetbrains.exposed.sql.select

class GitHubUsersDAO {
    private fun toGitHubUser(row: ResultRow) = GitHubUser(
        userID = row[GitHubUsers.user].value,
        username = row[GitHubUsers.username],
    )

    fun get(username: String) = GitHubUsers
        .select { GitHubUsers.username eq username }
        .map(::toGitHubUser)
        .first()

    fun create(userID: Int, username: String) = GitHubUsers
        .insert {
            it[GitHubUsers.user] = userID
            it[GitHubUsers.username] = username
        }.resultedValues!!.map(::toGitHubUser).first()
}

val gitHubUsersDAO = GitHubUsersDAO()
