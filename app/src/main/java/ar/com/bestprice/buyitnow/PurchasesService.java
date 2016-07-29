package ar.com.bestprice.buyitnow;

import com.google.gson.Gson;

import java.io.BufferedReader;
import java.io.BufferedWriter;
import java.io.InputStream;
import java.io.InputStreamReader;
import java.io.OutputStream;
import java.io.OutputStreamWriter;
import java.net.HttpURLConnection;
import java.net.URL;
import java.util.ArrayList;
import java.util.concurrent.Callable;
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

    public int savePurchases(ArrayList<Purchase> ps){

        final Purchases purchases = new Purchases();
        purchases.setPurchases(ps);

        Future<Integer>  task = service.submit(new Callable<Integer>() {

            @Override
            public Integer call() throws Exception {
                HttpURLConnection urlConnection = null;
                int responseCode = 203;

                URL url = new URL(Context.getContext().getServiceURL() + "/purchases");

                urlConnection = (HttpURLConnection) url.openConnection();
                urlConnection.setRequestMethod("POST");
                urlConnection.setRequestProperty("Authorization", Context.getContext().getSha1());
                urlConnection.setDoInput(true);
                urlConnection.setDoOutput(true);

                urlConnection.connect();

                Gson gson = new Gson();
                String it = gson.toJson(purchases);

                OutputStream os = urlConnection.getOutputStream();
                BufferedWriter writer = new BufferedWriter(new OutputStreamWriter(os, "UTF-8"));
                writer.write(it);
                writer.flush();
                writer.close();
                os.close();

                responseCode = urlConnection.getResponseCode();
                return responseCode;

        }});

        Integer status = 0;
        try {
            status = task.get();
        } catch (final InterruptedException | ExecutionException ex) {
            ex.printStackTrace();
        } finally {
            service.shutdownNow();
        }

        return status;
    }

    public int deletePurchase(Purchase ps){

        final Purchase p = ps;

        Future<Integer>  task = service.submit(new Callable<Integer>() {

            @Override
            public Integer call() throws Exception {
                HttpURLConnection urlConnection = null;
                int responseCode = 203;

                URL url = new URL(Context.getContext().getServiceURL() + "/purchases/" + p.getId());

                urlConnection = (HttpURLConnection) url.openConnection();

                urlConnection.setRequestMethod("DELETE");
                urlConnection.setRequestProperty("Authorization", Context.getContext().getSha1());
                urlConnection.setDoInput(true);
                urlConnection.setDoOutput(true);

                urlConnection.connect();


                OutputStream os = urlConnection.getOutputStream();
                BufferedWriter writer = new BufferedWriter(new OutputStreamWriter(os, "UTF-8"));
                //writer.write(it);
                writer.flush();
                writer.close();
                os.close();

                responseCode = urlConnection.getResponseCode();
                return responseCode;

            }});

        Integer status = 0;
        try {
            status = task.get();
        } catch (final InterruptedException | ExecutionException ex) {
            ex.printStackTrace();
        } finally {
            service.shutdownNow();
        }

        return status;
    }


    public String getPurchases() {

        Future<String>  task = service.submit(new Callable<String>() {
            @Override
            public String call() throws Exception{

                HttpURLConnection urlConnection = null;
                BufferedReader reader = null;
                String purchases = "";

                //URL url = new URL("http://10.33.117.120:8080/catalog/purchases?groupBy=month");
                URL url = new URL(Context.getContext().getServiceURL() + "/purchases?groupBy=month");

                // Create the request to OpenWeatherMap, and open the connection
                urlConnection = (HttpURLConnection) url.openConnection();
                urlConnection.setRequestMethod("GET");
                urlConnection.setRequestProperty("Authorization", Context.getContext().getSha1());
                urlConnection.connect();

                // Read the input stream into a String
                InputStream inputStream = urlConnection.getInputStream();
                StringBuffer buffer = new StringBuffer();
                if (inputStream == null) {
                    // Nothing to do.
                    return "";
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
                    return "";
                }

                purchases = buffer.toString();
                return purchases;
        }});

        String response = "";
        try {
            response =  task.get();
        } catch (InterruptedException e) {
            e.printStackTrace();
        } catch (ExecutionException e) {
            e.printStackTrace();
        }finally {
            service.shutdown();
        }

        return response;
    }

}
