package com.quentinguidee.plugins

import kotlin.io.path.Path

fun configureDataDirectory() {
    Path("data").toFile().mkdir()
}
