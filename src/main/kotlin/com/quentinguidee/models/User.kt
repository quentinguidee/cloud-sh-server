package com.quentinguidee.models

import kotlinx.serialization.Serializable
import org.jetbrains.exposed.sql.ResultRow
import org.jetbrains.exposed.sql.Table

object Users : Table() {
    val id = integer("id").autoIncrement()
    val username = varchar("username", 127).uniqueIndex()
    val name = varchar("name", 127)
    val email = varchar("email", 255)
    val profilePicture = varchar("profile_picture", 255)
    val role = varchar("role", 63).nullable()

    override val primaryKey = PrimaryKey(id)
}

@Serializable
data class User(
    val id: Int? = null,
    val username: String,
    val name: String,
    val email: String,
    val profilePicture: String,
    val role: String?,
)

fun Users.fromRow(row: ResultRow): User = User(
    id = row[id],
    username = row[username],
    name = row[name],
    email = row[email],
    profilePicture = row[profilePicture],
    role = row[role],
)
