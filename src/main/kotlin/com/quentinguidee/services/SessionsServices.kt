package com.quentinguidee.services

import com.quentinguidee.models.db.Session
import com.quentinguidee.models.db.Sessions
import org.jetbrains.exposed.sql.transactions.transaction

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
}

val sessionsServices = SessionsServices()
