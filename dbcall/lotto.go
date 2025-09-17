package dbcall

import (
    db "my-homepage/database"
    model "my-homepage/struct"
)

func GetLottoList() ([]model.GetLottoList, error) {
    rows, err := db.DB.Query(`
        SELECT index_no, no1, no2, no3, no4, no5, no6
        FROM number_list
        ORDER BY index_no DESC
    `)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var list []model.GetLottoList
    for rows.Next() {
        var lotto model.GetLottoList
        if err := rows.Scan(
            &lotto.IndexNo,
            &lotto.No1,
            &lotto.No2,
            &lotto.No3,
            &lotto.No4,
            &lotto.No5,
            &lotto.No6,
        ); err != nil {
            return nil, err
        }
        list = append(list, lotto)
    }
    return list, nil
}
