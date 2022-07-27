package com.quentinguidee.models.db

import kotlinx.serialization.json.JsonObject
import kotlinx.serialization.json.buildJsonObject
import kotlinx.serialization.json.put
import org.jetbrains.exposed.dao.IntEntity
import org.jetbrains.exposed.dao.IntEntityClass
import org.jetbrains.exposed.dao.id.EntityID
import org.jetbrains.exposed.dao.id.IntIdTable

object GitHubUsers : IntIdTable("github_users") {
    val user = reference("user", Users)
    val username = varchar("username", 255)
}

class GitHubUser(id: EntityID<Int>) : IntEntity(id) {
    companion object : IntEntityClass<GitHubUser>(GitHubUsers)

    var user by User referencedOn GitHubUsers.user
    var username by GitHubUsers.username

    fun toJSON(): JsonObject {
        return buildJsonObject {
            put("user", user.toJSON())
            put("username", username)
        }
    }
}
