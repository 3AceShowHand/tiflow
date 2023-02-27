package codec

import (
	"math/rand"

	"github.com/pingcap/tidb/parser/mysql"
	"github.com/pingcap/tidb/types"
	"github.com/pingcap/tidb/util/rowcodec"
	"github.com/pingcap/tiflow/cdc/model"
)

func NewLargeEvents(size int) *model.RowChangedEvent {
	value := make([]byte, size)
	rand.Read(value)

	return &model.RowChangedEvent{
		Table:     &model.TableName{Schema: "a", Table: "b"},
		TableInfo: &model.TableInfo{TableName: model.TableName{Schema: "a", Table: "b"}},
		Columns: []*model.Column{
			{
				Name:  "col-1",
				Type:  mysql.TypeBlob,
				Value: value,
			},
		},
		ColInfos: []rowcodec.ColInfo{{
			ID:            1000,
			IsPKHandle:    false,
			VirtualGenCol: false,
			Ft:            types.NewFieldType(mysql.TypeBlob),
		}},
	}
}
