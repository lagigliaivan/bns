package ar.com.bestprice.buyitnow.dto;

import java.util.List;

/**
 * Created by ivan on 07/04/16.
 */
public class PurchasesByMonthContainer {

    List<PurchasesByMonth> purchasesByMonth;

    public List<PurchasesByMonth> getPurchasesByMonth() {
        return purchasesByMonth;
    }

    public void setPurchasesByMonth(List<PurchasesByMonth> purchases) {
        if(purchases != null) {
            this.purchasesByMonth = purchases;
        }
    }

    public void addPurchasesByMonth(PurchasesByMonth purchasesByMonth) {

        if (this.purchasesByMonth != null) {
            this.purchasesByMonth.add(purchasesByMonth);
        }
    }
}
