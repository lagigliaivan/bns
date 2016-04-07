package ar.com.bestprice.buyitnow.dto;


import java.util.List;

/**
 * Created by ivan on 07/04/16.
 */
public class PurchasesByMonth {

    String month;
    List<Purchase> purchases;

    public String getMonth() {
        return month;
    }

    public void setMonth(String month) {
        this.month = month;
    }

    public List<Purchase> getPurchases() {
        return purchases;
    }

    public void setPurchases(List<Purchase> purchases) {
        this.purchases = purchases;
    }
}
