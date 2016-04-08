package ar.com.bestprice.buyitnow;

import android.content.Intent;
import android.os.Bundle;
import android.support.v7.app.AppCompatActivity;
import android.util.SparseArray;
import android.view.View;
import android.widget.ExpandableListView;

import com.google.gson.Gson;

import java.util.List;
import java.util.concurrent.ExecutionException;
import java.util.concurrent.ExecutorService;
import java.util.concurrent.Executors;
import java.util.concurrent.Future;

import ar.com.bestprice.buyitnow.dto.Item;
import ar.com.bestprice.buyitnow.dto.Purchase;
import ar.com.bestprice.buyitnow.dto.PurchasesByMonth;
import ar.com.bestprice.buyitnow.dto.PurchasesContainer;


public class MainActivity extends AppCompatActivity implements View.OnClickListener{

    private ExpandableListView listView = null;

    @Override
    protected void onCreate(Bundle savedInstanceState) {

        super.onCreate(savedInstanceState);
        setContentView(R.layout.activity_main);
        renderView();
    }

    private void renderView() {
        String jsonString = sendHttpRequest();
        PurchasesContainer purchasesContainer = parseJsonString(jsonString);

        ExpandableListView listView = getListView();
        MyExpandableListAdapter adapter = getListViewAdapter(purchasesContainer);
        listView.setAdapter(adapter);

    }

    private MyExpandableListAdapter getListViewAdapter(PurchasesContainer purchasesContainer) {
        SparseArray<Group> groups = createData(purchasesContainer.getPurchasesByMonth());
        MyExpandableListAdapter adapter = new MyExpandableListAdapter(this, groups);
        return adapter;

    }

    private String sendHttpRequest() {

        final ExecutorService service = Executors.newFixedThreadPool(1);
        final Future<String> task;
        String jsonString = "";

        task = service.submit(new GETServiceClient("http://10.33.117.120:8080/catalog/purchases?groupBy=month"));
        try {
            jsonString = task.get();
        } catch (final InterruptedException | ExecutionException ex) {
            ex.printStackTrace();
        } finally {
            service.shutdownNow();
        }

        return jsonString;
    }

    private ExpandableListView getListView() {

        if(this.listView == null) {
            this.listView = (ExpandableListView) findViewById(R.id.listView);
        }

        return listView;
    }

    public SparseArray<Group> createData( List<PurchasesByMonth> purchasesByMonth) {

        SparseArray<Group> groups = new SparseArray<>();
        int j = 0;
        for (PurchasesByMonth purchases:purchasesByMonth) {

            Group group = new Group(purchases.getMonth().toString());

            for (Purchase purchase : purchases.getPurchases()){

                for(Item item: purchase.getItems()) {
                    group.children.add(item);
                }
            }

            groups.append(j, group);
            j++;
        }

        return groups;

    }
    private PurchasesContainer parseJsonString(String json){

        Gson gson = new Gson();
        PurchasesContainer p = gson.fromJson(json, PurchasesContainer.class);
        return p;
    }

    //private static final int ADD_NEW_PURCHASE = 9001;

    public void onClick(View view) {
        if(view.getId() == R.id.add_new_purchase) {
            Intent intent = new Intent(this.getApplicationContext(), AddNewPurchaseActivity.class);
            startActivity(intent);
        }else{
            String jsonString = sendHttpRequest();
            PurchasesContainer purchasesContainer = parseJsonString(jsonString);
            ExpandableListView listView = getListView();
            MyExpandableListAdapter adapter = getListViewAdapter(purchasesContainer);
            listView.setAdapter(adapter);
        }
    }


    @Override
    protected void onStart() {
        super.onStart();
        // The activity is about to become visible.
    }
    @Override
    protected void onResume() {
        super.onResume();
        // The activity has become visible (it is now "resumed").
    }
    @Override
    protected void onPause() {
        super.onPause();
        // Another activity is taking focus (this activity is about to be "paused").
    }
    @Override
    protected void onStop() {
        super.onStop();
        // The activity is no longer visible (it is now "stopped")
    }
    @Override
    protected void onDestroy() {
        super.onDestroy();
        // The activity is about to be destroyed.
    }

}
