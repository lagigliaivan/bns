package ar.com.bestprice.buyitnow.dto;

import java.util.List;

/**
 * Created by ivan on 07/04/16.
 */
public class PurchasesContainer {

    List<PurchasesByMonth> purchasesByMonth;

    public List<PurchasesByMonth> getPurchasesByMonth() {
        return purchasesByMonth;
    }

    public void setPurchasesByMonth(List<PurchasesByMonth> purchasesByMonth) {
        this.purchasesByMonth = purchasesByMonth;
    }
}
