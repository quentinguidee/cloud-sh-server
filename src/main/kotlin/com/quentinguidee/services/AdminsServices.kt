package com.quentinguidee.services

import com.quentinguidee.plugins.DB_CONFIG_PATH
import com.quentinguidee.plugins.DatabaseConfig
import com.quentinguidee.plugins.resetDatabase
import kotlinx.serialization.json.Json
import kotlinx.serialization.json.encodeToJsonElement
import java.io.File
import kotlin.io.path.Path

class AdminsServices {
    fun reset() {
        resetDatabase()
        File(Path("data", "buckets").toString()).deleteRecursively()
    }

    fun saveDatabaseConfig(config: DatabaseConfig) {
        val json = Json.encodeToJsonElement(config)
        DB_CONFIG_PATH.toFile().writeText(json.toString())
    }
}

val adminsServices = AdminsServices()
