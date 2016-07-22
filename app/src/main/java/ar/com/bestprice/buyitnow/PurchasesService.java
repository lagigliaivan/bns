package ar.com.bestprice.buyitnow;

import java.util.ArrayList;
import java.util.concurrent.ExecutionException;
import java.util.concurrent.ExecutorService;
import java.util.concurrent.Executors;
import java.util.concurrent.Future;

import ar.com.bestprice.buyitnow.dto.Purchase;
import ar.com.bestprice.buyitnow.dto.Purchases;

/**
 * Created by elagiglia on 22/7/16.
 */
public class PurchasesService {

    final ExecutorService service = Executors.newFixedThreadPool(1);


    public void savePurchases(ArrayList<Purchase> ps){
        Purchases purchases = new Purchases();
        purchases.setPurchases(ps);
        String serviceURL = Context.getContext().getServiceURL();
        Future<Integer>  task = service.submit(new POSTServiceClient(serviceURL + "/purchases", purchases));

        try {
            Integer status = task.get();
        } catch (final InterruptedException | ExecutionException ex) {
            ex.printStackTrace();
        } finally {
            service.shutdownNow();
        }

    }




}
