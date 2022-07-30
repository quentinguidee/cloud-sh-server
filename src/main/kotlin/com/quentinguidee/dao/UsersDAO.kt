package com.quentinguidee.dao

import com.quentinguidee.models.User
import com.quentinguidee.models.Users
import org.jetbrains.exposed.sql.ResultRow
import org.jetbrains.exposed.sql.insert
import org.jetbrains.exposed.sql.select

class UsersDAO {
    private fun toUser(row: ResultRow) = User(
        id = row[Users.id].value,
        username = row[Users.username],
        name = row[Users.name],
        email = row[Users.email],
        profilePicture = row[Users.profilePicture],
        role = row[Users.role],
    )

    fun get(userID: Int) = Users
        .select { Users.id eq userID }
        .map(::toUser)
        .first()

    fun get(username: String) = Users
        .select { Users.username eq username }
        .map(::toUser)
        .first()

    fun create(username: String, name: String, email: String, profilePicture: String, role: String? = null) =
        Users.insert {
            it[Users.username] = username
            it[Users.name] = name
            it[Users.email] = email
            it[Users.profilePicture] = profilePicture
            it[Users.role] = role
        }.resultedValues?.map(::toUser)!!.first()
}

val usersDAO = UsersDAO()
