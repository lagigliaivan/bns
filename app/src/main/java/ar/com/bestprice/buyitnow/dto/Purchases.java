package ar.com.bestprice.buyitnow.dto;

import java.util.List;

/**
 * Created by ivan on 08/04/16.
 */
public class Purchases {
    List<Purchase> purchases;

    public List<Purchase> getPurchases() {
        return purchases;
    }

    public void setPurchases(List<Purchase> purchases) {
        this.purchases = purchases;
    }
}
