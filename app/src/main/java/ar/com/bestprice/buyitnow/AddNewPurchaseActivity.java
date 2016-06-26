package ar.com.bestprice.buyitnow;

import android.content.Intent;
import android.os.Bundle;
import android.support.v7.app.AppCompatActivity;
import android.support.v7.widget.Toolbar;
import android.util.Log;
import android.view.Menu;
import android.view.MenuItem;
import android.widget.ArrayAdapter;
import android.widget.ListView;

import com.google.android.gms.common.api.CommonStatusCodes;
import com.google.android.gms.vision.barcode.Barcode;

import java.text.SimpleDateFormat;
import java.util.ArrayList;
import java.util.Date;
import java.util.List;
import java.util.Locale;
import java.util.TimeZone;
import java.util.concurrent.ExecutionException;
import java.util.concurrent.ExecutorService;
import java.util.concurrent.Executors;
import java.util.concurrent.Future;

import ar.com.bestprice.buyitnow.barcodereader.BarcodeCaptureActivity;
import ar.com.bestprice.buyitnow.dto.Item;
import ar.com.bestprice.buyitnow.dto.Purchase;
import ar.com.bestprice.buyitnow.dto.Purchases;

/**
 * Created by ivan on 08/04/16.
 */
public class AddNewPurchaseActivity extends AppCompatActivity{

    ListView listView = null;
    ArrayAdapter<String> adapter = null;
    List<Item> items = new ArrayList<>();

    @Override
    protected void onCreate(Bundle savedInstanceState) {

        super.onCreate(savedInstanceState);
        setContentView(R.layout.activity_add_new_purchase_tool_bar);

        listView = (ListView) findViewById(R.id.listview_show_items_in_a_purchase);

        Item item = (Item)getIntent().getSerializableExtra("Item");

        items.add(item);

        List itemsAsString = new ArrayList();

        for (Item i: items ){
            itemsAsString.add(i.toString());
        }

        // Define a new Adapter
        // First parameter - Context
        // Second parameter - Layout for the row
        // Third parameter - ID of the TextView to which the data is written
        // Forth - the Array of data

        adapter = new ArrayAdapter<>(this,
                android.R.layout.simple_list_item_1, android.R.id.text1, itemsAsString);

        // Assign adapter to ListView
        listView.setAdapter(adapter);


        Toolbar toolbar = (Toolbar) findViewById(R.id.new_purchase_toolbar);

        setSupportActionBar(toolbar);
        getSupportActionBar().setDisplayShowTitleEnabled(false);
        getSupportActionBar().setDisplayHomeAsUpEnabled(false);
        getSupportActionBar().setDisplayShowTitleEnabled(false);

    }

    private static final int RC_BARCODE_CAPTURE = 9001;


    private static String TAG = "BarCode Reader";
    private static final int ADD_ITEM = 9002;

    @Override
    protected void onActivityResult(int requestCode, int resultCode, Intent data) {

        if (requestCode == RC_BARCODE_CAPTURE) {

            if (resultCode == CommonStatusCodes.SUCCESS) {
                if (data != null) {
                    Barcode barcode = data.getParcelableExtra(BarcodeCaptureActivity.BarcodeObject);

                    Log.d(TAG, "Barcode read: " + barcode.displayValue);

                    Intent intent = new Intent(this.getApplicationContext(), AddItemActivity.class);
                    intent.putExtra("BarCode", barcode.displayValue);
                    startActivityForResult(intent, ADD_ITEM);

                } else {
                    Log.d(TAG, "No barcode captured, intent data is null");
                }

            }
        }
        else{
            if (resultCode == CommonStatusCodes.SUCCESS && data != null) {

                    Item item = (Item) data.getSerializableExtra("Item");

                    items.add(item);

                    List itemsAsString = new ArrayList();

                    for (Item i : items) {
                        itemsAsString.add(i.toString());
                    }

                    ArrayAdapter adapter = new ArrayAdapter<>(this,
                            android.R.layout.simple_list_item_1, android.R.id.text1, itemsAsString);

                    listView.setAdapter(adapter);
            }
        } /*else {
            super.onActivityResult(requestCode, resultCode, data);
        }*/
    }

    @Override
    public boolean onCreateOptionsMenu(Menu menu) {
        getMenuInflater().inflate(R.menu.save_purchase_activity_toolbar_menu, menu);
        return true;
    }

    @Override
    public boolean onOptionsItemSelected(MenuItem item) {

        // Handle action bar item clicks here. The action bar will
        // automatically handle clicks on the Home/Up button, so long
        // as you specify a parent activity in AndroidManifest.xml.

        switch (item.getItemId()){

            case R.id.add_item_no_barcode:
                Intent intent = new Intent(this.getApplicationContext(), AddItemActivity.class);
                intent.putExtra(Constants.CALLING_ACTIVITY, Constants.NEW_PURCHASE);
                startActivityForResult(intent, Constants.NEW_PURCHASE);
                break;

            case R.id.add_item_barcode:

                startActivity(new Intent(this.getApplicationContext(), BarcodeCaptureActivity.class));
                break;

            case R.id.save_purchase:

                final ExecutorService service = Executors.newFixedThreadPool(1);
                final Future<Integer> task;


                //2016-05-05T18:54:03.5102707-03:00
                SimpleDateFormat datetime = new SimpleDateFormat ("yyyy-MM-dd'T'HH:mm:ss.SSSSSSSZZZZZ", Locale.US);
                datetime.setTimeZone(TimeZone.getTimeZone("UTC"));

                Date date = new Date(System.currentTimeMillis());

                Purchases purchases = new Purchases();

                Purchase purchase = new Purchase();
                purchase.setItems(items);
                purchase.setTime(datetime.format(date));


                ArrayList<Purchase> ps = new ArrayList<>();
                ps.add(purchase);

                purchases.setPurchases(ps);


                //task = service.submit(new POSTServiceClient("http://10.33.117.120:8080/catalog/purchases", purchases));
                String serviceURL = Context.getContext().getServiceURL();
                task = service.submit(new POSTServiceClient(serviceURL + "/purchases", purchases));


                try {
                    Integer status = task.get();
                } catch (final InterruptedException | ExecutionException ex) {
                    ex.printStackTrace();
                } finally {
                    service.shutdownNow();
                }


                finish();
        }

        return super.onOptionsItemSelected(item);
    }
}
