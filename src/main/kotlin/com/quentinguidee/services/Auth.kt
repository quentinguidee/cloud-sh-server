package com.quentinguidee.services

import com.quentinguidee.client
import com.quentinguidee.models.*
import io.ktor.client.call.*
import io.ktor.client.request.*
import io.ktor.http.*
import kotlinx.serialization.SerialName
import kotlinx.serialization.Serializable
import org.jetbrains.exposed.dao.load
import org.jetbrains.exposed.sql.transactions.transaction
import java.util.*

@Serializable
data class GitHubUserBody(
    val email: String,
    val name: String,
    @SerialName("avatar_url")
    val avatarURL: String,
    val login: String,
)

class AuthService {
    suspend fun githubUser(username: String) = transaction {
        GitHubUser
            .find { GitHubUsers.username eq username }
            .firstOrNull()
    }

    suspend fun fetchGitHubUser(token: String): GitHubUserBody {
        return client
            .request("https://api.github.com/user") {
                headers {
                    append(HttpHeaders.Authorization, "token $token")
                }
            }
            .body()
    }

    private suspend fun createAccount(
        username: String,
        name: String,
        email: String,
        profilePicture: String
    ): Session = transaction {
        val user = User.new {
            this.username = username
            this.name = name
            this.email = email
            this.profilePicture = profilePicture
        }

        GitHubUser.new {
            this.user = user
            this.username = username
        }

        return@transaction Session.new {
            this.user = user
            this.token = UUID.randomUUID().toString()
        }
    }

    suspend fun createAccount(gitHubUser: GitHubUserBody) = authService.createAccount(
        username = gitHubUser.login,
        name = gitHubUser.name,
        email = gitHubUser.email,
        profilePicture = gitHubUser.avatarURL,
    )

    suspend fun getAccount(username: String) = transaction {
        val user = User
            .find { Users.username eq username }
            .first()

        return@transaction Session
            .find { Sessions.user eq user.id }
            .first()
            .load(Session::user)
    }
}

val authService = AuthService()
