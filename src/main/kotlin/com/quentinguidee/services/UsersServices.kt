package com.quentinguidee.services

import com.quentinguidee.dao.usersDAO
import org.jetbrains.exposed.sql.transactions.transaction

class UsersServices {
    fun get(username: String) = transaction {
        usersDAO.get(username)
    }

    fun get(userID: Int) = transaction {
        usersDAO.get(userID)
    }

    fun updateInfo(userID: Int, name: String?, email: String?, profilePicture: String?) = transaction {
        usersDAO.update(
            userID = userID,
            name = name,
            email = email,
            profilePicture = profilePicture
        )
        usersDAO.get(userID)
    }
}

val usersServices = UsersServices()
