package com.quentinguidee.services

import com.quentinguidee.models.db.Session
import com.quentinguidee.models.db.Sessions
import com.quentinguidee.models.db.User
import org.jetbrains.exposed.sql.transactions.transaction
import java.util.*

class SessionsServices {
    suspend fun session(token: String) = transaction {
        Session
            .find { Sessions.token eq token }
            .first()
    }

    suspend fun revokeSession(token: String) = transaction {
        Session
            .find { Sessions.token eq token }
            .first()
            .delete()
    }

    suspend fun createSession(user: User) = transaction {
        Session.new {
            this.user = user
            this.token = UUID.randomUUID().toString()
        }
    }
}

val sessionsServices = SessionsServices()
