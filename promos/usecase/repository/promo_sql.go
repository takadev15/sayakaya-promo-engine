package repository

const (
    queryInsertBirthdayPromo = `
        INSERT INTO promos (
            promo_code, 
            user_id, 
            promo_type_id, 
            valid_until,
            status
        )
        VALUES ($1, $2, $3, $4, $5)
    `

    querySelectPromoTypes = `
        SELECT
            id,
            name,
            rule
        FROM promo_type
    `

    queryCheckPromoUserActive = `
        SELECT
            id
        FROM promos
        WHERE user_id = $1 AND
        status = 0 OR status = 1
    `
)
