package ar.com.bestprice.buyitnow;

import android.content.Intent;
import android.os.Bundle;
import android.support.v7.app.AppCompatActivity;
import android.support.v7.widget.Toolbar;
import android.util.SparseArray;
import android.view.Menu;
import android.view.MenuItem;
import android.widget.ExpandableListView;

import com.google.gson.Gson;

import java.util.List;
import java.util.concurrent.ExecutionException;
import java.util.concurrent.ExecutorService;
import java.util.concurrent.Executors;
import java.util.concurrent.Future;

import ar.com.bestprice.buyitnow.barcodereader.BarcodeCaptureActivity;
import ar.com.bestprice.buyitnow.dto.Item;
import ar.com.bestprice.buyitnow.dto.Purchase;
import ar.com.bestprice.buyitnow.dto.PurchasesByMonth;
import ar.com.bestprice.buyitnow.dto.PurchasesContainer;


public class MainActivity extends AppCompatActivity {


    private ExpandableListView listView = null;

    @Override
    protected void onCreate(Bundle savedInstanceState) {

        super.onCreate(savedInstanceState);
        setContentView(R.layout.activity_main_tool_bar);
        renderView();
    }

    private void renderView() {
        String jsonString = sendHttpRequest();

        if(jsonString != null) {

            PurchasesContainer purchasesContainer = parseJsonString(jsonString);
            ExpandableListView listView = getListView();

            MyExpandableListAdapter adapter = getListViewAdapter(purchasesContainer);
            listView.setAdapter(adapter);

            Toolbar toolbar = (Toolbar) findViewById(R.id.main_tool_bar);

            setSupportActionBar(toolbar);
            getSupportActionBar().setDisplayShowTitleEnabled(false);
            getSupportActionBar().setDisplayHomeAsUpEnabled(false);
            getSupportActionBar().setDisplayShowTitleEnabled(false);
        }

    }

    private MyExpandableListAdapter getListViewAdapter(PurchasesContainer purchasesContainer) {

        SparseArray<Group> groups = createData(purchasesContainer.getPurchasesByMonth());
        return new MyExpandableListAdapter(this, groups);

    }

    private String sendHttpRequest() {

        final ExecutorService service = Executors.newFixedThreadPool(1);
        final Future<String> task;
        String jsonString = "";

        //task = service.submit(new GETServiceClient("http://10.33.117.120:8080/catalog/purchases?groupBy=month"));
        String serviceURL = Context.getContext().getServiceURL();
        String user = Context.getContext().getUser();

        task = service.submit(new GETServiceClient(serviceURL + "/purchases?groupBy=month", Context.getContext().getSha1()));

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
            this.listView = (ExpandableListView) findViewById(R.id.listView_show_purchases);
        }

        return listView;
    }

    private SparseArray<Group> createData( List<PurchasesByMonth> purchasesByMonth) {

        SparseArray<Group> groups = new SparseArray<>();
        int j = 0;
        for (PurchasesByMonth purchases:purchasesByMonth) {

            Group group = new Group(purchases.getMonth());

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

    @Override
    public boolean onCreateOptionsMenu(Menu menu) {
        getMenuInflater().inflate(R.menu.main_activity_toolbar_menu, menu);
        return true;
    }

    @Override
    public boolean onOptionsItemSelected(MenuItem item) {

        // Handle action bar item clicks here. The action bar will
        // automatically handle clicks on the Home/Up button, so long
        // as you specify a parent activity in AndroidManifest.xml.

        switch (item.getItemId()){

            case R.id.add_item_no_barcode:

                startActivity(new Intent(this.getApplicationContext(), AddItemActivity.class));
                break;

            case R.id.add_item_barcode:

                startActivity(new Intent(this.getApplicationContext(), BarcodeCaptureActivity.class));
                break;

            case R.id.refresh_purchases:

                String jsonString = sendHttpRequest();
                PurchasesContainer purchasesContainer = parseJsonString(jsonString);
                MyExpandableListAdapter adapter = getListViewAdapter(purchasesContainer);

                ExpandableListView listView = getListView();
                listView.setAdapter(adapter);
                break;
        }

        return super.onOptionsItemSelected(item);
    }
}
