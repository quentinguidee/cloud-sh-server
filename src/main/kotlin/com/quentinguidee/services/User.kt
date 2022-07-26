package com.quentinguidee.services

import com.quentinguidee.models.User
import com.quentinguidee.models.Users
import com.quentinguidee.models.fromRow
import org.jetbrains.exposed.sql.select
import org.jetbrains.exposed.sql.transactions.transaction

class UserService {
    fun getUser(username: String): User? {
        return transaction {
            val query = Users
                .select { Users.username eq username }
                .singleOrNull() ?: return@transaction null

            query.let { row -> return@let Users.fromRow(row) }
        }
    }

}

val userService = UserService()
