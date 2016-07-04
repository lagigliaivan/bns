package ar.com.bestprice.buyitnow;
import java.util.ArrayList;
import java.util.List;

import ar.com.bestprice.buyitnow.dto.Item;
import ar.com.bestprice.buyitnow.dto.Purchase;

public class PurchasesGroup {

    public Month getMonth() {
        return month;
    }

    private Month month;

    public final List<Item> children = new ArrayList<>();
    private float purchasesTotalPrice = 0;

    public String getString() {
        return month.toString();
    }

    public PurchasesGroup(Month month) {
        this.month = month;
    }


    public float getPurchasesTotalPrice() {
        return purchasesTotalPrice;
    }

    public void setPurchasesTotalPrice(float purchasesTotalPrice) {
        this.purchasesTotalPrice = purchasesTotalPrice;
    }
}
