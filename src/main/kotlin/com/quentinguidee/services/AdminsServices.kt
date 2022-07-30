package com.quentinguidee.services

import com.quentinguidee.plugins.initDatabase
import com.quentinguidee.plugins.resetDatabase

class AdminsServices {
    suspend fun reset() {
        resetDatabase()
        initDatabase()
    }
}

val adminsServices = AdminsServices()
