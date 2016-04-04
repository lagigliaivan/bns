package ar.com.bestprice.buyitnow;

import android.os.Bundle;
import android.support.v7.app.AppCompatActivity;
import android.util.Log;
import android.view.View;
import android.widget.EditText;

import com.google.gson.Gson;

import java.io.BufferedReader;
import java.io.BufferedWriter;
import java.io.IOException;
import java.io.InputStream;
import java.io.InputStreamReader;
import java.io.OutputStream;
import java.io.OutputStreamWriter;
import java.net.HttpURLConnection;
import java.net.URL;
import java.net.URLEncoder;
import java.util.ArrayList;
import java.util.List;

/**
 * Created by ivan on 01/04/16.
 */
public class AddItemActivity extends AppCompatActivity implements View.OnClickListener{

    @Override
    protected void onCreate(Bundle savedInstanceState) {

        super.onCreate(savedInstanceState);

        setContentView(R.layout.activity_add_item);

        String barcodeId = getIntent().getStringExtra("BarCode");

        EditText t = (EditText)findViewById(R.id.editText);

        t.setText(barcodeId);
    }

    @Override
    public void onClick(View v) {
        HttpURLConnection urlConnection = null;
        BufferedReader reader = null;

        EditText id = (EditText)findViewById(R.id.editText);
        EditText description = (EditText)findViewById(R.id.editText2);
        EditText price = (EditText)findViewById(R.id.editText3);

        // Will contain the raw JSON response as a string.

        try {
            // http://openweathermap.org/API#forecast
            //URL url = new URL("http://10.33.117.120:8080/catalog/products/");
            URL url = new URL("http://192.168.0.7:8080/catalog/products/");

            urlConnection = (HttpURLConnection) url.openConnection();
            urlConnection.setRequestMethod("POST");
            urlConnection.setDoInput(true);
            urlConnection.setDoOutput(true);

            urlConnection.connect();

            Item item = new Item();
            item.setId(id.getText().toString());
            item.setDescription(description.getText().toString());
            item.setPrice(Float.valueOf(price.getText().toString()));

            List<Item> list = new ArrayList<>();
            list.add(item);

            Items items = new Items();
            items.setItems(list);

            Gson gson = new Gson();

            String it = gson.toJson(items);

            OutputStream os = urlConnection.getOutputStream();
            BufferedWriter writer = new BufferedWriter(
                    new OutputStreamWriter(os, "UTF-8"));
            //writer.write(getPostDataString(postDataParams));
            writer.write(it);
            writer.flush();
            writer.close();
            os.close();


            int responseCode = urlConnection.getResponseCode();

            finish();

        } catch (IOException e) {
            Log.e("PlaceholderFragment", "Error ", e);
            // If the code didn't successfully get the weather data, there's no point in attemping
            // to parse it.

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
}
