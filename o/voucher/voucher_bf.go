package voucher

import (
	"ehelp/x/rest"
	"errors"
	"strings"
	"time"
)

func (v *Voucher) ValidateVoucher() (*Voucher, error) {
	var vou, err = GetVoucherByID(v.Code)
	if err != nil {
		return nil, err
	}
	return vou.Validate(0)
}

func (vou *Voucher) Validate(typeSer int) (*Voucher, error) {
	if typeSer != 0 && vou.ServiceType != typeSer {
		return nil, rest.BadRequestValid(errors.New("Mã " + vou.Code + " không sử dụng cho loại hình này"))
	}
	var timeNow = time.Now().Unix()
	if vou.StartTime > timeNow || vou.EndTime < timeNow {
		return nil, rest.BadRequestValid(errors.New("Mã đã hết hạn sử dụng!"))
	}
	if !vou.Active && !vou.AutoActive {
		return nil, rest.BadRequestValid(errors.New("Mã chưa được kích hoạt!"))
	}
	if vou.Quantity != 0 && vou.Count >= vou.Quantity {
		return nil, rest.BadRequestValid(errors.New("Mã quá số lượng cho phép!"))
	}

	return vou, nil
}

func GetVoucherByID(vouCode string) (vouRes *Voucher, err error) {
	if len(vouCode) == 0 {
		err = rest.BadRequestValid(errors.New("Nhập mã khuyến mại!"))
		return
	}
	vouCode = strings.ToUpper(vouCode)
	if VoucherCache != nil && len(VoucherCache) > 0 {
		for _, vou := range VoucherCache {
			if strings.ToUpper(vou.Code) == vouCode {
				vouRes = vou
				return
			}
		}
		if vouRes == nil {
			err = rest.BadRequestValid(errors.New("Mã không tồn tại!"))
		}
	} else {
		vouRes, err = GetVoucherCode(vouCode)
		if err != nil {
			if rest.IsNotFound(err) {
				return nil, rest.BadRequestValid(errors.New("Mã không tồn tại!"))
			}
			return nil, rest.BadRequestValid(errors.New("Có lỗi xảy ra. Vui lòng thử lại!"))
		}
		if vouRes != nil {
			var timeNow = time.Now().Unix()
			if vouRes.StartTime > timeNow || vouRes.EndTime < timeNow {
				return nil, rest.BadRequestValid(errors.New("Mã đã hết hạn!"))
			}
		}
	}
	return
}
