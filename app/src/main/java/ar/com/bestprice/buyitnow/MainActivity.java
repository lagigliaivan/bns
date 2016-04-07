package ar.com.bestprice.buyitnow;
import android.content.Intent;
import android.os.Bundle;

import android.support.design.widget.FloatingActionButton;
import android.support.v7.app.AppCompatActivity;
import android.util.Log;
import android.util.SparseArray;
import android.view.View;
import android.widget.ExpandableListView;
import android.widget.FrameLayout;


import com.google.android.gms.common.api.CommonStatusCodes;
import com.google.android.gms.vision.barcode.Barcode;
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


public class MainActivity extends AppCompatActivity implements View.OnClickListener{

    private ExpandableListView listView = null;

    private static final String TAG = "BarcodeMain";

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

        task = service.submit(new ServiceClient("http://10.33.117.120:8080/catalog/purchases?groupBy=month"));
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

            FrameLayout footerLayout = (FrameLayout) getLayoutInflater().inflate(R.layout.footer, null);
            FrameLayout headerLayout = (FrameLayout) getLayoutInflater().inflate(R.layout.header, null);

            FloatingActionButton btnAddItem = (FloatingActionButton) footerLayout.findViewById(R.id.add_item_button);
            btnAddItem.setOnClickListener(this);

            FloatingActionButton btnRefresh = (FloatingActionButton) headerLayout.findViewById(R.id.refresh_button);
            btnRefresh.setOnClickListener(this);


            listView.addHeaderView(headerLayout);
            listView.addFooterView(footerLayout);
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

    private static final int RC_BARCODE_CAPTURE = 9001;

    public void onClick(View view) {
        if(view.getId() == R.id.add_item_button) {
            Intent intent = new Intent(this.getApplicationContext(), BarcodeCaptureActivity.class);
            intent.putExtra(BarcodeCaptureActivity.AutoFocus, true);
            intent.putExtra(BarcodeCaptureActivity.UseFlash, false);
            startActivityForResult(intent, RC_BARCODE_CAPTURE);
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
    /**
     * Called when an activity you launched exits, giving you the requestCode
     * you started it with, the resultCode it returned, and any additional
     * data from it.  The <var>resultCode</var> will be
     * {@link #RESULT_CANCELED} if the activity explicitly returned that,
     * didn't return any result, or crashed during its operation.
     * <p/>
     * <p>You will receive this call immediately before onResume() when your
     * activity is re-starting.
     * <p/>
     *
     * @param requestCode The integer request code originally supplied to
     *                    startActivityForResult(), allowing you to identify who this
     *                    result came from.
     * @param resultCode  The integer result code returned by the child activity
     *                    through its setResult().
     * @param data        An Intent, which can return result data to the caller
     *                    (various data can be attached to Intent "extras").
     * @see #startActivityForResult
     * @see #createPendingResult
     * @see #setResult(int)
     */
    @Override
    protected void onActivityResult(int requestCode, int resultCode, Intent data) {
        if (requestCode == RC_BARCODE_CAPTURE) {
            if (resultCode == CommonStatusCodes.SUCCESS) {
                if (data != null) {
                    Barcode barcode = data.getParcelableExtra(BarcodeCaptureActivity.BarcodeObject);

                    //statusMessage.setText(R.string.barcode_success);
                    //barcodeValue.setText(barcode.displayValue);
                    Log.d(TAG, "Barcode read: " + barcode.displayValue);

                    Intent intent = new Intent(this.getApplicationContext(), AddItemActivity.class);
                    intent.putExtra("BarCode", barcode.displayValue);
                    startActivity(intent);
                    //startActivityForResult(intent, 12);
                } else {
                    //statusMessage.setText(R.string.barcode_failure);
                    Log.d(TAG, "No barcode captured, intent data is null");
                }
            } else {
                Log.d(TAG, "Exiting from CRUD");
            }
        }
        else {
            super.onActivityResult(requestCode, resultCode, data);
        }
    }

}
