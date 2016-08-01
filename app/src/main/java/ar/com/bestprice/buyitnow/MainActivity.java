package ar.com.bestprice.buyitnow;

import android.app.AlertDialog;
import android.content.DialogInterface;
import android.content.Intent;
import android.os.Bundle;
import android.support.annotation.NonNull;
import android.support.v7.app.AppCompatActivity;
import android.support.v7.widget.Toolbar;
import android.text.InputType;
import android.view.Menu;
import android.view.MenuItem;
import android.view.View;
import android.widget.EditText;
import android.widget.ExpandableListView;
import android.widget.TextView;

import com.google.gson.Gson;

import java.util.ArrayList;
import java.util.HashMap;
import java.util.List;
import java.util.Map;

import ar.com.bestprice.buyitnow.dto.Item;
import ar.com.bestprice.buyitnow.dto.Purchase;
import ar.com.bestprice.buyitnow.dto.PurchasesByMonth;
import ar.com.bestprice.buyitnow.dto.PurchasesByMonthContainer;


public class MainActivity extends AppCompatActivity {


    private ExpandableListView listView = null;

    //Contains the purchases returned by the server
    private PurchasesByMonthContainer purchasesContainer = null;

    @Override
    protected void onCreate(Bundle savedInstanceState) {

        super.onCreate(savedInstanceState);
        setContentView(R.layout.activity_main_tool_bar);
        renderView();
    }

    private void renderView() {

        renderPurchasesList();

        Toolbar toolbar = (Toolbar) findViewById(R.id.main_tool_bar);

        setSupportActionBar(toolbar);

        getSupportActionBar().setDisplayHomeAsUpEnabled(false);
        getSupportActionBar().setDisplayShowTitleEnabled(false);

    }

    private MyExpandableListAdapter getListViewAdapter(PurchasesByMonthContainer purchasesContainer) {

        Map<Integer, PurchasesGroup> groups = getSortedPurchasesGroups(purchasesContainer.getPurchasesByMonth());
        return new MyExpandableListAdapter(this, groups);

    }

    private String sendHttpRequest() {

        String jsonString = null;

        try {
            PurchasesService service = new PurchasesService(Context.getContext());
            jsonString = service.getPurchases();
        } catch (Exception ex) {
            ex.printStackTrace();
        }

        return jsonString;
    }

    private ExpandableListView getListView() {

        if(this.listView == null) {
            this.listView = (ExpandableListView) findViewById(R.id.listView_show_purchases);
        }

        return listView;
    }

    private Map<Integer, PurchasesGroup> getSortedPurchasesGroups(List<PurchasesByMonth> purchasesByMonth) {

        Map<Month, PurchasesByMonth> sortedPurchases = sortPurchases(purchasesByMonth);

        Map<Integer, PurchasesGroup> groups = new HashMap<>();

        int j = 0;
        for (Month month : Month.values()){

            if (sortedPurchases.get(month) != null){

                PurchasesGroup purchasesGroup = new PurchasesGroup(month);

                for (Purchase purchase : sortedPurchases.get(month).getPurchases()){

                    for(Item item: purchase.getItems()) {
                        item.setTime(purchase.getTime());
                        purchasesGroup.addItem(item);
                    }
                    purchasesGroup.addPurchase(purchase);
                }
                groups.put(j, purchasesGroup);
                j++;
            }

        }

        return groups;

    }

    @NonNull
    private Map<Month, PurchasesByMonth> sortPurchases(List<PurchasesByMonth> purchasesByMonth) {

        Map<Month, PurchasesByMonth> sortedPurchases = new HashMap<>();
        for (PurchasesByMonth purchases : purchasesByMonth) {
            sortedPurchases.put(Month.valueOf(purchases.getMonth().toUpperCase()), purchases);
        }

        return sortedPurchases;
    }

    @Override
    public boolean onOptionsItemSelected(MenuItem item) {

        // Handle action bar item clicks here. The action bar will
        // automatically handle clicks on the Home/Up button, so long
        // as you specify a parent activity in AndroidManifest.xml.

        switch (item.getItemId()){

            case R.id.add_item:

                Intent intent = new Intent(this.getApplicationContext(), AddItemActivity.class);
                intent.putExtra(Constants.CALLING_ACTIVITY, Constants.MAIN_ACTIVITY);
                startActivity(intent);
                break;

            case R.id.refresh_purchases:

                renderPurchasesList();
                break;

            case R.id.search_purchases:

                AlertDialog.Builder builder = new AlertDialog.Builder(this);
                builder.setTitle("Ingrese descripcion a buscar");

                // Set up the input
                final EditText input = new EditText(this);
                // Specify the type of input expected; this, for example, sets the input as a password, and will mask the text
                input.setInputType(InputType.TYPE_CLASS_TEXT);
                builder.setView(input);

                // Set up the buttons
                builder.setPositiveButton("OK", new DialogInterface.OnClickListener() {

                    @Override
                    public void onClick(DialogInterface dialog, int which) {

                        String m_Text = input.getText().toString();

                        List purchasesByMonth = searchStringInPurchases(m_Text, purchasesContainer.getPurchasesByMonth());

                        PurchasesByMonthContainer container = new PurchasesByMonthContainer();
                        container.setPurchasesByMonth(purchasesByMonth);

                        renderList(container);
                    }
                });
                builder.setNegativeButton("Cancel", new DialogInterface.OnClickListener() {
                    @Override
                    public void onClick(DialogInterface dialog, int which) {
                        dialog.cancel();
                    }
                });

                builder.show();

                break;
        }

        return super.onOptionsItemSelected(item);
    }

    private void renderPurchasesList() {

        String jsonString = sendHttpRequest();


//        String jsonString =  "{\"purchasesByMonth\":[\n" +
//                "\t\t{ \"month\" : \"January\",\n" +
//                "\t\t  \"purchases\":[\n" +
//                "\t\t\t\t{\n" +
//                "\t\t  \t\t\"time\":\"2016-04-12T00:06:22.364Z\",\n" +
//                "\t\t  \t\t\"items\":[\n" +
//                "\t\t\t   \t\t{\n" +
//                "\t\t\t     \t\t\"id\":\"1233123\",\n" +
//                "\t\t\t     \t\t\"description\":\"no se\",\n" +
//                "\t\t\t     \t\t\"price\":12.0,\n" +
//                "\t\t\t     \t\t\"category\":\"SALUD\"\n" +
//                "\t\t\t   \t\t},\n" +
//                "\t\t\t\t\t{\n" +
//                "\t\t\t     \t\t\"id\":\"1234343123\",\n" +
//                "\t\t\t     \t\t\"description\":\"very long long long description of a product which has really a long long long description. Not sure what will happen\",\n" +
//                "\t\t\t     \t\t\"price\":12.0,\n" +
//                "\t\t\t     \t\t\"category\":\"SALUD\"\n" +
//                "\t\t\t   \t\t}\n" +
//                "\n" +
//                "\t\t\t  \t\t]\t\n" +
//                "\t\t\t\t}\n" +
//                "\t  \t\t ]\n" +
//                "\t\t},\n" +
//                "\t\t{ \"month\" : \"March\",\n" +
//                "\t\t  \"purchases\":[\n" +
//                "\t\t\t\t{\n" +
//                "\t\t  \t\t\"time\":\"2016-04-12T00:06:22.364Z\",\n" +
//                "\t\t  \t\t\"items\":[\n" +
//                "\t\t\t   \t\t{\n" +
//                "\t\t\t     \t\t\"id\":\"1233123\",\n" +
//                "\t\t\t     \t\t\"description\":\"no se\",\n" +
//                "\t\t\t     \t\t\"price\":32.0,\n" +
//                "\t\t\t     \t\t\"category\":\"SALUD\"\n" +
//                "\t\t\t   \t\t},\n" +
//                "\t\t\t   \t\t{\n" +
//                "                        \"id\":\"1233123\",\n" +
//                "                        \"description\":\"no se\",\n" +
//                "                        \"price\":32.0,\n" +
//                "                        \"category\":\"MERCADERIA\"\n" +
//                "                    },\n" +
//                "\t\t\t\t\t{\n" +
//                "\t\t\t     \t\t\"id\":\"1234343123\",\n" +
//                "\t\t\t     \t\t\"description\":\"very long long long description of a product which has really a long long long description. Not sure what will happen\",\n" +
//                "\t\t\t     \t\t\"price\":22.0,\n" +
//                "\t\t\t     \t\t\"category\":\"SALUD\"\n" +
//                "\t\t\t   \t\t}\n" +
//                "\n" +
//                "\t\t\t  \t\t]\t\n" +
//                "\t\t\t\t}\n" +
//                "\t  \t\t ]\n" +
//                "\t\t},{ \"month\" : \"April\",\n" +
//                "          \t\t  \"purchases\":[\n" +
//                "          \t\t\t\t{\n" +
//                "          \t\t  \t\t\"time\":\"2016-04-12T00:06:22.364Z\",\n" +
//                "          \t\t  \t\t\"items\":[\n" +
//                "          \t\t\t   \t\t{\n" +
//                "          \t\t\t     \t\t\"id\":\"1233123\",\n" +
//                "          \t\t\t     \t\t\"description\":\"no se\",\n" +
//                "          \t\t\t     \t\t\"price\":2.0,\n" +
//                "          \t\t\t     \t\t\"category\":\"DIVERSION\"\n" +
//                "          \t\t\t   \t\t},\n" +
//                "          \t\t\t   \t\t{\n" +
//                "                                  \"id\":\"1233123\",\n" +
//                "                                  \"description\":\"no se\",\n" +
//                "                                  \"price\":3.0,\n" +
//                "                                  \"category\":\"MERCADERIA\"\n" +
//                "                              },\n" +
//                "          \t\t\t\t\t{\n" +
//                "          \t\t\t     \t\t\"id\":\"1234343123\",\n" +
//                "          \t\t\t     \t\t\"description\":\"very long long long description of a product which has really a long long long description. Not sure what will happen\",\n" +
//                "          \t\t\t     \t\t\"price\":22.0,\n" +
//                "          \t\t\t     \t\t\"category\":\"SERVICIOS\"\n" +
//                "          \t\t\t   \t\t}\n" +
//                "\n" +
//                "          \t\t\t  \t\t]\n" +
//                "          \t\t\t\t}\n" +
//                "          \t  \t\t ]\n" +
//                "          \t\t}\n" +
//                "\n" +
//                "]}";


        if (jsonString != null && !jsonString.isEmpty()){
            purchasesContainer = parseJsonString(jsonString);
            renderList(purchasesContainer);
        }
    }

    private void renderList(PurchasesByMonthContainer purchasesContainer) {

        ExpandableListView listView = getListView();
        MyExpandableListAdapter adapter = getListViewAdapter(purchasesContainer);
        adapter.setParent(listView);
        listView.setAdapter(adapter);

        listView.setOnGroupClickListener(new ExpandableListView.OnGroupClickListener() {
            @Override
            public boolean onGroupClick(ExpandableListView parent, View v, int groupPosition, long id) {
                return false;
            }
        });

        Map<Integer, PurchasesGroup> purchases = getSortedPurchasesGroups(purchasesContainer.getPurchasesByMonth());

        float accumPurchases = 0;

        for (PurchasesGroup group : purchases.values()) {
            accumPurchases += group.getPurchasesTotalPrice();
        }

        float purchasesAverage = 0;
        if (purchases.size()!= 0) {
            purchasesAverage = accumPurchases / purchases.size();
        }
        TextView average = (TextView) findViewById(R.id.average);
        TextView accumulated = (TextView) findViewById(R.id.accumulated);

        average.setText(String.format("Promedio mensual: $%.2f", purchasesAverage));
        accumulated.setText(String.format("Acumulado: $%.2f", accumPurchases));

    }

    private PurchasesByMonthContainer parseJsonString(String json){

        Gson gson = new Gson();
        PurchasesByMonthContainer p = gson.fromJson(json, PurchasesByMonthContainer.class);
        return p;
    }

    @Override
    public boolean onCreateOptionsMenu(Menu menu) {
        getMenuInflater().inflate(R.menu.main_activity_toolbar_menu, menu);
        return true;
    }

    @Override
    protected void onResume() {
        super.onResume();
        renderPurchasesList();
    }

    private List<PurchasesByMonth> searchStringInPurchases(String pattern, List<PurchasesByMonth> purchasesByMonths){

        List<PurchasesByMonth> byMonths = new ArrayList<>();


        for(PurchasesByMonth pByMonth : purchasesByMonths){

            PurchasesByMonth pByMonthWhereItemWasFound = null;

            for (Purchase purchase : pByMonth.getPurchases()){


                Purchase pWhereItemWasFound = null;

                for(Item item: purchase.getItems()){

                    if(item.getDescription().contains(pattern)){

                        if (pWhereItemWasFound == null){
                            pWhereItemWasFound = new Purchase();
                            pWhereItemWasFound.setTime(purchase.getTime());
                        }

                        pWhereItemWasFound.addItem(item);
                    }
                }

                if(pWhereItemWasFound != null){
                    if(pByMonthWhereItemWasFound == null) {
                        pByMonthWhereItemWasFound = new PurchasesByMonth();
                        pByMonthWhereItemWasFound.setMonth(pByMonth.getMonth());
                    }
                    pByMonthWhereItemWasFound.addPurchase(pWhereItemWasFound);
                }
            }

            if(pByMonthWhereItemWasFound != null) {
                byMonths.add(pByMonthWhereItemWasFound);
            }
        }

        return byMonths;
    }
}
