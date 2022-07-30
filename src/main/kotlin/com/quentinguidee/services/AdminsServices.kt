package com.quentinguidee.services

import com.quentinguidee.plugins.initDatabase
import com.quentinguidee.plugins.resetDatabase
import java.io.File
import kotlin.io.path.Path

class AdminsServices {
    fun reset() {
        resetDatabase()
        initDatabase()

        File(Path("data", "buckets").toString()).deleteRecursively()
    }
}

val adminsServices = AdminsServices()
