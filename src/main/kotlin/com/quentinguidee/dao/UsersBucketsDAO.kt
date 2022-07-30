package com.quentinguidee.dao

import com.quentinguidee.models.AccessType
import com.quentinguidee.models.UserBucket
import com.quentinguidee.models.UsersBuckets
import org.jetbrains.exposed.sql.ResultRow
import org.jetbrains.exposed.sql.and
import org.jetbrains.exposed.sql.insert
import org.jetbrains.exposed.sql.select
import java.util.*

class UsersBucketsDAO {
    private fun toUserBucket(row: ResultRow) = UserBucket(
        bucketUUID = row[UsersBuckets.bucket].value.toString(),
        userID = row[UsersBuckets.user].value,
        accessType = row[UsersBuckets.accessType],
    )

    fun get(bucketUUID: UUID, userID: Int) = UsersBuckets
        .select {
            UsersBuckets.user eq userID and
                    (UsersBuckets.bucket eq bucketUUID)
        }
        .map(::toUserBucket)
        .first()

    fun create(bucketUUID: UUID, userID: Int, accessType: AccessType) = UsersBuckets
        .insert {
            it[UsersBuckets.bucket] = bucketUUID
            it[UsersBuckets.user] = userID
            it[UsersBuckets.accessType] = accessType
        }.resultedValues!!.map(::toUserBucket).first()
}

val usersBucketsDAO = UsersBucketsDAO()
