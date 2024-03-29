package com.quentinguidee.plugins

import io.ktor.client.request.*
import io.ktor.client.statement.*
import io.ktor.http.*
import io.ktor.server.testing.*
import kotlin.test.Test
import kotlin.test.assertEquals

class RoutingTest {
    @Test
    fun testGetPing() = testApplication {
        application {
            configureRouting()
        }

        client.get("/ping").apply {
            assertEquals(HttpStatusCode.OK, status)
            assertEquals("pong", bodyAsText())
        }
    }
}
