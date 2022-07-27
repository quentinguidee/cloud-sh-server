package com.quentinguidee.models.db

import kotlinx.serialization.json.buildJsonObject
import kotlinx.serialization.json.put
import org.jetbrains.exposed.dao.IntEntity
import org.jetbrains.exposed.dao.IntEntityClass
import org.jetbrains.exposed.dao.id.EntityID
import org.jetbrains.exposed.dao.id.IntIdTable
import org.jetbrains.exposed.sql.transactions.transaction

object Users : IntIdTable() {
    val username = varchar("username", 127).uniqueIndex()
    val name = varchar("name", 127)
    val email = varchar("email", 255)
    val profilePicture = varchar("profile_picture", 255)
    val role = varchar("role", 63).nullable()
}

class User(id: EntityID<Int>) : IntEntity(id) {
    companion object : IntEntityClass<User>(Users)

    var username by Users.username
    var name by Users.name
    var email by Users.email
    var profilePicture by Users.profilePicture
    var role by Users.role

    var buckets by Bucket via UserBuckets

    fun toJSON() = transaction {
        return@transaction buildJsonObject {
            put("id", id)
            put("username", username)
            put("name", name)
            put("email", email)
            put("profile_picture", profilePicture)
            put("role", role)
        }
    }
}
