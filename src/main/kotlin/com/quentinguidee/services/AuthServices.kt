package com.quentinguidee.services

import com.quentinguidee.client
import com.quentinguidee.dao.oAuthMethodsDAO
import com.quentinguidee.dao.oAuthUsersDAO
import com.quentinguidee.dao.sessionsDAO
import com.quentinguidee.dao.usersDAO
import com.quentinguidee.models.OAuthMethodPrivate
import com.quentinguidee.models.Session
import io.ktor.client.call.*
import io.ktor.client.request.*
import io.ktor.http.*
import kotlinx.serialization.SerialName
import kotlinx.serialization.Serializable
import org.jetbrains.exposed.sql.transactions.transaction

@Serializable
data class GitHubUserBody(
    val email: String,
    var name: String,
    @SerialName("avatar_url")
    val avatarURL: String,
    val login: String,
)

class AuthServices {
    fun oAuthUser(username: String, method: OAuthMethodPrivate) = transaction {
        oAuthUsersDAO.get(username, method)
    }

    suspend fun fetchGitHubUser(token: String): GitHubUserBody {
        val githubUser: GitHubUserBody = client
            .request("https://api.github.com/user") {
                headers {
                    append(HttpHeaders.Authorization, "token $token")
                }
            }
            .body()

        if (githubUser.name.isBlank())
            githubUser.name = githubUser.login

        return githubUser
    }

    private fun createAccount(
        method: OAuthMethodPrivate,
        username: String,
        name: String,
        email: String,
        profilePicture: String
    ): Session =
        transaction {
            val user = usersDAO.create(
                username = username,
                name = name,
                email = email,
                profilePicture = profilePicture,
                role = "admin"
            )

            oAuthUsersDAO.create(user.id, username, method)

            return@transaction sessionsDAO.create(user.id)
        }

    fun createAccount(gitHubUser: GitHubUserBody, method: OAuthMethodPrivate) = authServices.createAccount(
        method = method,
        username = gitHubUser.login,
        name = gitHubUser.name,
        email = gitHubUser.email,
        profilePicture = gitHubUser.avatarURL,
    )

    fun session(username: String) = transaction {
        val user = usersDAO.get(username)
        return@transaction sessionsDAO.get(user.id)
    }

    fun methods() = transaction {
        oAuthMethodsDAO.getAll()
    }

    fun methodsPrivate() = transaction {
        oAuthMethodsDAO.getAllPrivate()
    }

    fun method(name: String) = transaction {
        oAuthMethodsDAO.get(name)
    }

    fun methodPrivate(name: String) = transaction {
        oAuthMethodsDAO.getPrivate(name)
    }

    fun createMethod(
        name: String,
        displayName: String,
        color: String,
        clientID: String,
        clientSecret: String,
        authorizeURL: String,
        accessTokenURL: String,
        redirectURL: String,
    ) = transaction {
        oAuthMethodsDAO.create(
            name,
            displayName,
            color,
            clientID,
            clientSecret,
            authorizeURL,
            accessTokenURL,
            redirectURL,
        )
    }
}

val authServices = AuthServices()
