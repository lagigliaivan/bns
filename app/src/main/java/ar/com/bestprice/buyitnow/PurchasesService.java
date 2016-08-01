package ar.com.bestprice.buyitnow;

import android.support.annotation.NonNull;

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
import java.util.ArrayList;
import java.util.concurrent.Callable;
import java.util.concurrent.ExecutorService;
import java.util.concurrent.Executors;
import java.util.concurrent.Future;

import ar.com.bestprice.buyitnow.dto.Purchase;
import ar.com.bestprice.buyitnow.dto.Purchases;

/**
 * Created by elagiglia on 22/7/16.
 */
public class PurchasesService {

    Context context;

    public PurchasesService(Context ctx){
        context = ctx;
    }

    final ExecutorService service = Executors.newFixedThreadPool(1);


    public int savePurchases(@NonNull ArrayList<Purchase> ps) throws Exception{

        final Purchases purchases = new Purchases();
        purchases.setPurchases(ps);

        Future<Integer>  task = service.submit(new Callable<Integer>() {

            @Override
            public Integer call() throws Exception {

                HttpURLConnection urlConnection = getHttpURLConnection("POST", "/purchases");

                OutputStream os = urlConnection.getOutputStream();
                Gson gson = new Gson();
                String it = gson.toJson(purchases);
                BufferedWriter writer = new BufferedWriter(new OutputStreamWriter(os, "UTF-8"));
                writer.write(it);
                writer.flush();
                writer.close();
                os.close();

                int responseCode = urlConnection.getResponseCode();

                return responseCode;
        }});

        int status = task.get();
        service.shutdownNow();

        return status;
    }

    @NonNull
    public int deletePurchase(@NonNull Purchase ps) throws Exception{

        final Purchase p = ps;

        Future<Integer>  task = service.submit(new Callable<Integer>() {

            @Override
            public Integer call() throws Exception {

                HttpURLConnection urlConnection = getHttpURLConnection("DELETE", "/purchases/" + p.getId());
                OutputStream os = urlConnection.getOutputStream();

                BufferedWriter writer = new BufferedWriter(new OutputStreamWriter(os, "UTF-8"));
                writer.flush();
                writer.close();
                os.close();

                return  urlConnection.getResponseCode();
            }});

        int status = task.get();
        service.shutdownNow();

        return status;
    }

    @NonNull
    public String getPurchases() throws Exception{

        Future<String>  task = service.submit(new Callable<String>() {
            @Override
            public String call() throws Exception{

                HttpURLConnection urlConnection = getHttpURLConnection("GET", "/purchases?groupBy=month");

                InputStream inputStream = urlConnection.getInputStream();
                StringBuffer buffer = new StringBuffer();
                if (inputStream == null) {
                    return "";
                }

                BufferedReader reader = new BufferedReader(new InputStreamReader(inputStream));

                String line;
                while ((line = reader.readLine()) != null) {
                    buffer.append(line + "\n");
                }

                if (buffer.length() == 0) {
                    return "";
                }

                return buffer.toString();

        }});

        String response =  task.get();
        service.shutdown();

        return response;
    }

    @NonNull
    private HttpURLConnection getHttpURLConnection(@NonNull String httpMethod, @NonNull String resource) throws IOException {

        HttpURLConnection urlConnection;

        URL url = new URL(context.getServiceURL() + resource);
        urlConnection = (HttpURLConnection) url.openConnection();
        urlConnection.setRequestMethod(httpMethod);
        urlConnection.setRequestProperty("Authorization", context.getSha1());
        urlConnection.setDoInput(true);
        urlConnection.setDoOutput(true);
        urlConnection.connect();

        return urlConnection;
    }
}
