package ar.com.bestprice.buyitnow.dto;


import java.util.ArrayList;
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

        if(purchases != null) {
            this.purchases = purchases;
        }
    }

    public void addPurchase(Purchase purchase) {

        if(getPurchases() == null){
           setPurchases(new ArrayList<Purchase>());
        }

        getPurchases().add(purchase);
    }
}
