package com.quentinguidee.services

import com.quentinguidee.models.db.Session
import com.quentinguidee.models.db.Sessions
import org.jetbrains.exposed.sql.transactions.transaction

class SessionServices {
    suspend fun session(token: String) = transaction {
        Session
            .find { Sessions.token eq token }
            .first()
    }
}

val sessionServices = SessionServices()
