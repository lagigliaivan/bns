package ar.com.bestprice.buyitnow;

import android.app.Activity;
import android.content.Intent;
import android.os.Bundle;
import android.os.StrictMode;
import android.support.design.widget.FloatingActionButton;
import android.support.v7.app.AppCompatActivity;
import android.util.Log;
import android.util.SparseArray;
import android.view.View;
import android.widget.ArrayAdapter;
import android.widget.Button;
import android.widget.ExpandableListView;
import android.widget.FrameLayout;
import android.widget.ListView;
import android.widget.TextView;

import com.google.android.gms.common.api.CommonStatusCodes;
import com.google.android.gms.vision.barcode.Barcode;

import java.io.BufferedReader;
import java.io.IOException;
import java.io.InputStream;
import java.io.InputStreamReader;
import java.net.HttpURLConnection;
import java.net.URL;
import java.util.ArrayList;
import java.util.Arrays;
import java.util.List;
import com.google.gson.Gson;


import ar.com.bestprice.buyitnow.barcodereader.BarcodeCaptureActivity;


public class MainActivity extends AppCompatActivity implements View.OnClickListener{

    private TextView statusMessage;
    private TextView barcodeValue;
    private static final String TAG = "BarcodeMain";

    @Override
    protected void onCreate(Bundle savedInstanceState) {
        StrictMode.ThreadPolicy policy = new StrictMode.ThreadPolicy.Builder().permitAll().build();

        StrictMode.setThreadPolicy(policy);

        super.onCreate(savedInstanceState);

        String[] data = {
                "Mon 6/23 - Sunny - 31/17",
                "Tue 6/24 - Foggy - 21/8",
                "Wed 6/25 - Cloudy - 22/17",
                "Thurs 6/26 - Rainy - 18/11",
                "Fri 6/27 - Foggy - 21/10",
                "Sat 6/28 - TRAPPED IN WEATHERSTATION - 23/18",
                "Sun 6/29 - Sunny - 20/7"
        };

        //List<String> items = new ArrayList<>(Arrays.asList(data));
        List<Item> items = getAllItems();
        //getAllItems();

        setContentView(R.layout.activity_main);

        createData(items);
        ExpandableListView listView = (ExpandableListView) findViewById(R.id.listView);
        //ListView view =  (ListView) findViewById(R.id.item_list_view);
        FrameLayout footerLayout = (FrameLayout) getLayoutInflater().inflate(R.layout.footer, null);
        FloatingActionButton btnPostYourEnquiry = (FloatingActionButton) footerLayout.findViewById(R.id.add_item_button);
        btnPostYourEnquiry.setOnClickListener(this);
        listView.addFooterView(footerLayout);

        MyExpandableListAdapter adapter = new MyExpandableListAdapter(this, groups);
        listView.setAdapter(adapter);

       /* ArrayAdapter<String> mForecastAdapter =
                new ArrayAdapter<>(
                        view.getContext(), // The current context (this activity)
                        R.layout.activity_main, // The name of the layout ID.
                        R.id.item_text_view, // The ID of the textview to populate.
                        items);

        view.setAdapter(mForecastAdapter);*/

    }

    SparseArray<Group> groups = new SparseArray<>();
    public void createData( List<Item> items) {

        for (int j = 0; j < 5; j++) {

            Group group = new Group("Fecha " + j);
            for (int i = 0; i < items.size(); i++) {
                group.children.add(items.get(i));
            }
            groups.append(j, group);
        }
    }

    private List getAllItems(){

        // These two need to be declared outside the try/catch
        // so that they can be closed in the finally block.
        HttpURLConnection urlConnection = null;
        BufferedReader reader = null;

        // Will contain the raw JSON response as a string.
        String items = null;

        try {
            // http://openweathermap.org/API#forecast
            URL url = new URL("http://10.33.117.120:8080/catalog/products/");

            // Create the request to OpenWeatherMap, and open the connection
            urlConnection = (HttpURLConnection) url.openConnection();
            urlConnection.setRequestMethod("GET");
            urlConnection.connect();

            // Read the input stream into a String
            InputStream inputStream = urlConnection.getInputStream();
            StringBuffer buffer = new StringBuffer();
            if (inputStream == null) {
                // Nothing to do.
                return new ArrayList();
            }
            reader = new BufferedReader(new InputStreamReader(inputStream));

            String line;
            while ((line = reader.readLine()) != null) {
                // Since it's JSON, adding a newline isn't necessary (it won't affect parsing)
                // But it does make debugging a *lot* easier if you print out the completed
                // buffer for debugging.
                buffer.append(line + "\n");
            }

            if (buffer.length() == 0) {
                // Stream was empty.  No point in parsing.
                return new ArrayList();
            }


            items = buffer.toString();
            Gson gson = new Gson();

            ArrayList it = gson.fromJson(items, ArrayList.class);
            return it;

        } catch (IOException e) {
            Log.e("PlaceholderFragment", "Error ", e);
            // If the code didn't successfully get the weather data, there's no point in attemping
            // to parse it.
            return new ArrayList();
        } finally{
            if (urlConnection != null) {
                urlConnection.disconnect();
            }
            if (reader != null) {
                try {
                    reader.close();
                } catch (final IOException e) {
                    Log.e("PlaceholderFragment", "Error closing stream", e);
                }
            }
        }


    }


    private static final int RC_BARCODE_CAPTURE = 9001;
    public void onClick(View view) {
        if(view.getId() == R.id.add_item_button) {
            Intent intent = new Intent(this.getApplicationContext(), BarcodeCaptureActivity.class);
            intent.putExtra(BarcodeCaptureActivity.AutoFocus, true);
            intent.putExtra(BarcodeCaptureActivity.UseFlash, false);
            startActivityForResult(intent, RC_BARCODE_CAPTURE);
        }else{
            System.out.println(((TextView) view).getText());
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
                } else {
                    //statusMessage.setText(R.string.barcode_failure);
                    Log.d(TAG, "No barcode captured, intent data is null");
                }
            } else {
                //statusMessage.setText(String.format(getString(R.string.barcode_error),
                //        CommonStatusCodes.getStatusCodeString(resultCode)));
            }
        }
        else {
            super.onActivityResult(requestCode, resultCode, data);
        }
    }

}
