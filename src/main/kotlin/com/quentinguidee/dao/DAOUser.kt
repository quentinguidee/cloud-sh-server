package com.quentinguidee.dao

import com.quentinguidee.models.User
import com.quentinguidee.models.Users

class DAOUser {
    fun get(username: String) = User
        .find { Users.username eq username }
        .singleOrNull()
}

var daoUser = DAOUser()
